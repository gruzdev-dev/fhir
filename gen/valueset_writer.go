package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gruzdev-dev/fhir/tools/text"
)

type ValueSetConstants struct {
	TypeName string
	URL      string
	Codes    []Code
}

func (g *Generator) GenerateValueSets() error {
	requiredURLs, err := g.ExtractRequiredValueSets()
	if err != nil {
		return fmt.Errorf("extract required value sets: %w", err)
	}

	valueSets, codeSystems, err := g.LoadValueSets()
	if err != nil {
		return fmt.Errorf("load value sets: %w", err)
	}

	filteredValueSets := make([]ValueSetResource, 0)
	for url, vs := range valueSets {
		if requiredURLs[url] {
			filteredValueSets = append(filteredValueSets, vs)
		}
	}

	constants := make([]ValueSetConstants, 0)
	for _, vs := range filteredValueSets {
		codes, err := ExtractCodesFromValueSet(vs, codeSystems)
		if err != nil {
			continue
		}
		if len(codes) == 0 {
			continue
		}

		codes = removeDuplicateCodes(codes)

		typeName := normalizeTypeName(vs.Name)
		if typeName == "" {
			continue
		}

		constants = append(constants, ValueSetConstants{
			TypeName: typeName,
			URL:      vs.URL,
			Codes:    codes,
		})
	}

	sort.Slice(constants, func(i, j int) bool {
		return constants[i].TypeName < constants[j].TypeName
	})

	return g.writeConstantsFiles(constants)
}

func removeDuplicateCodes(codes []Code) []Code {
	seen := make(map[string]bool)
	result := make([]Code, 0)
	for _, code := range codes {
		if !seen[code.Code] {
			seen[code.Code] = true
			result = append(result, code)
		}
	}
	return result
}

func normalizeTypeName(name string) string {
	if name == "" {
		return ""
	}

	typeName := text.TitleCase(name)

	var builder strings.Builder
	for _, r := range typeName {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			builder.WriteRune(r)
		}
	}
	typeName = builder.String()

	if !text.IsValidGoIdentifier(typeName) {
		if len(typeName) > 0 && typeName[0] >= '0' && typeName[0] <= '9' {
			typeName = "ValueSet" + typeName
		}
		if !text.IsValidGoIdentifier(typeName) {
			return ""
		}
	}

	return typeName
}

func normalizeConstantName(typeName, code string) string {
	parts := strings.FieldsFunc(code, func(r rune) bool {
		return r == '-' || r == '_' || r == ' ' || r == '.'
	})

	var builder strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		part = text.TitleCase(part)
		builder.WriteString(part)
	}
	constantName := builder.String()

	var cleanBuilder strings.Builder
	for _, r := range constantName {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			cleanBuilder.WriteRune(r)
		}
	}
	constantName = cleanBuilder.String()

	if constantName == "" {
		constantName = "Code"
	} else if constantName[0] >= '0' && constantName[0] <= '9' {
		constantName = "Code" + constantName
	}

	return typeName + constantName
}

func (g *Generator) writeConstantsFiles(constants []ValueSetConstants) error {
	groups := groupByAlphabet(constants)

	for rangeName, groupConstants := range groups {
		filename := fmt.Sprintf("codes_%s.go", rangeName)
		if err := g.writeConstantsFile(filename, groupConstants); err != nil {
			return fmt.Errorf("write file %s: %w", filename, err)
		}
	}
	return nil
}

func groupByAlphabet(constants []ValueSetConstants) map[string][]ValueSetConstants {
	groups := make(map[string][]ValueSetConstants)

	for _, c := range constants {
		if len(c.TypeName) == 0 {
			continue
		}
		firstChar := strings.ToLower(string(c.TypeName[0]))
		var rangeName string
		switch {
		case firstChar >= "a" && firstChar <= "f":
			rangeName = "a_f"
		case firstChar >= "g" && firstChar <= "l":
			rangeName = "g_l"
		case firstChar >= "m" && firstChar <= "r":
			rangeName = "m_r"
		case firstChar >= "s" && firstChar <= "z":
			rangeName = "s_z"
		default:
			rangeName = "other"
		}
		groups[rangeName] = append(groups[rangeName], c)
	}

	return groups
}

func (g *Generator) writeConstantsFile(filename string, constants []ValueSetConstants) error {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "package models\n\n")

	usedConstantNames := make(map[string]bool)
	usedTypeNames := make(map[string]bool)
	for _, c := range constants {
		if usedTypeNames[c.TypeName] {
			continue
		}
		usedTypeNames[c.TypeName] = true

		if err := g.writeValueSetConstants(&buf, c, usedConstantNames); err != nil {
			return err
		}
		fmt.Fprintf(&buf, "\n")
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		_ = os.WriteFile("debug_valueset_failed.go", buf.Bytes(), 0644)
		return fmt.Errorf("format error for %s: %w. Check debug_valueset_failed.go", filename, err)
	}

	filePath := filepath.Join(g.OutputPath, filename)
	return os.WriteFile(filePath, formatted, 0644)
}

func (g *Generator) writeValueSetConstants(buf *bytes.Buffer, c ValueSetConstants, usedConstantNames map[string]bool) error {
	fmt.Fprintf(buf, "// %s represents codes from %s\n", c.TypeName, c.URL)
	fmt.Fprintf(buf, "type %s string\n\n", c.TypeName)
	fmt.Fprintf(buf, "const (\n")

	for _, code := range c.Codes {
		constantName := normalizeConstantName(c.TypeName, code.Code)

		originalName := constantName
		counter := 2
		for usedConstantNames[constantName] {
			constantName = fmt.Sprintf("%s_%d", originalName, counter)
			counter++
		}
		usedConstantNames[constantName] = true

		codeValue := strings.ReplaceAll(code.Code, "`", "\\`")
		codeValue = strings.ReplaceAll(codeValue, "\"", "\\\"")

		fmt.Fprintf(buf, "\t%s %s = %q\n", constantName, c.TypeName, codeValue)
	}

	fmt.Fprintf(buf, ")\n")
	return nil
}
