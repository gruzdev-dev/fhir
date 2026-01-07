package gen

import (
	"fmt"
	"strings"

	"github.com/gruzdev-dev/fhir/tools/text"
)

type FieldInfo struct {
	Name       string
	GoType     string
	JSONTag    string
	BSONTag    string
	Comment    string
	Min        int
	MaxLength  *int
	Pattern    string
	Fixed      any
	IsRequired bool
	Path       string
}

func (g *Generator) ProcessElements(name string, elements []ElementDefinition, def StructureDefinition) map[string][]FieldInfo {
	structs := make(map[string][]FieldInfo)

	isPrimitiveType := def.Kind == "primitive-type"
	originalName := def.Name

	for _, el := range elements {
		if isPrimitiveType && el.Path == originalName {
		} else if el.Path == name || strings.Contains(el.ID, "._") {
			continue
		}

		if isPrimitiveType && el.Path == originalName+".value" {
			var valueType string
			switch originalName {
			case "string", "uri", "code", "id", "markdown", "canonical", "url", "base64Binary", "oid", "xhtml",
				"date", "dateTime", "instant", "time", "uuid":
				valueType = "*string"
			case "boolean":
				valueType = "*bool"
			case "integer", "unsignedInt", "positiveInt":
				valueType = "*int"
			case "integer64":
				valueType = "*int64"
			case "decimal":
				valueType = "*float64"
			default:
				continue
			}

			jsonTag := "`json:\"value\"`"
			if el.Min == 0 {
				jsonTag = "`json:\"value,omitempty\"`"
			}
			bsonTag := generateBSONTag("Value", el.Min == 0)

			structs[name] = append(structs[name], FieldInfo{
				Name:       "Value",
				GoType:     valueType,
				JSONTag:    jsonTag,
				BSONTag:    bsonTag,
				Comment:    el.Short,
				Min:        el.Min,
				MaxLength:  el.MaxLength,
				Pattern:    el.Pattern,
				Fixed:      el.Fixed,
				IsRequired: el.Min > 0,
				Path:       el.Path,
			})
			continue
		}

		if isPrimitiveType && el.Path == originalName+".id" {
			jsonTag := "`json:\"id\"`"
			if el.Min == 0 {
				jsonTag = "`json:\"id,omitempty\"`"
			}
			bsonTag := generateBSONTag("Id", el.Min == 0)

			structs[name] = append(structs[name], FieldInfo{
				Name:       "Id",
				GoType:     "*string",
				JSONTag:    jsonTag,
				BSONTag:    bsonTag,
				Comment:    el.Short,
				Min:        el.Min,
				MaxLength:  el.MaxLength,
				Pattern:    el.Pattern,
				Fixed:      el.Fixed,
				IsRequired: el.Min > 0,
				Path:       el.Path,
			})
			continue
		}

		parts := strings.Split(el.Path, ".")
		lastPart := parts[len(parts)-1]

		if isPrimitiveType {
			if lastPart == "id" || lastPart == "value" {
				continue
			}
			lastPartLower := strings.ToLower(lastPart)
			if needsFHIRPrefix(lastPart) || needsFHIRPrefix(lastPartLower) {
				continue
			}
			if def, ok := g.Definitions[lastPart]; ok && def.Kind == "primitive-type" {
				continue
			}
			if def, ok := g.Definitions[lastPartLower]; ok && def.Kind == "primitive-type" {
				continue
			}
		}

		if strings.HasPrefix(lastPart, "_") || lastPart == "extension" || lastPart == "modifierExtension" {
			continue
		}

		var structName string
		if len(parts) == 1 {
			structName = name
		} else {
			parentPath := strings.Join(parts[:len(parts)-1], ".")
			structName = g.deriveNestedTypeName(parentPath)
		}

		if el.ContentReference != "" && strings.HasPrefix(el.ContentReference, "#") {
			refPath := strings.TrimPrefix(el.ContentReference, "#")
			goType := g.deriveNestedTypeName(refPath)
			if goType == "" {
				goType = "any"
			}

			cleanName := strings.ReplaceAll(lastPart, "[x]", "")
			cleanName = strings.ReplaceAll(cleanName, "-", "_")
			cleanName = text.TitleCase(cleanName)

			if cleanName == "" {
				continue
			}

			if el.Max == "*" {
				goType = "[]" + goType
			} else if el.Min == 0 && !strings.HasPrefix(goType, "[]") && goType != "json.RawMessage" && goType != "any" {
				goType = "*" + goType
			} else if el.Min > 0 && !strings.HasPrefix(goType, "[]") {
				isPrimitive := goType == "bool" || goType == "string" || goType == "int" || goType == "int64" || goType == "float64" || goType == "json.RawMessage" || goType == "any"
				if !isPrimitive {
					goType = "*" + goType
				}
			}

			escapedTag := strings.ReplaceAll(lastPart, "\"", "\\\"")
			jsonTag := fmt.Sprintf("`json:\"%s\"`", escapedTag)
			if el.Min == 0 {
				jsonTag = fmt.Sprintf("`json:\"%s,omitempty\"`", escapedTag)
			}
			bsonTag := generateBSONTag(cleanName, el.Min == 0)

			structs[structName] = append(structs[structName], FieldInfo{
				Name:       cleanName,
				GoType:     goType,
				JSONTag:    jsonTag,
				BSONTag:    bsonTag,
				Comment:    el.Short,
				Min:        el.Min,
				MaxLength:  el.MaxLength,
				Pattern:    el.Pattern,
				Fixed:      el.Fixed,
				IsRequired: el.Min > 0,
				Path:       el.Path,
			})
			continue
		}

		if strings.Contains(lastPart, "[x]") {
			baseName := strings.ReplaceAll(lastPart, "[x]", "")
			for _, fhirType := range el.Type {
				singleTypeEl := el
				singleTypeEl.Type = []ElementDataType{fhirType}

				goType := g.mapGoType(singleTypeEl)

				typeSuffix := text.TitleCase(fhirType.Code)
				fieldName := text.TitleCase(baseName) + typeSuffix
				jsonTag := fmt.Sprintf("`json:\"%s%s\"`", baseName, typeSuffix)
				if el.Min == 0 {
					jsonTag = fmt.Sprintf("`json:\"%s%s,omitempty\"`", baseName, typeSuffix)
				}
				bsonTag := generateBSONTag(fieldName, el.Min == 0)

				if !strings.HasPrefix(goType, "*") && !strings.HasPrefix(goType, "[]") && goType != "any" {
					goType = "*" + goType
				}

				structs[structName] = append(structs[structName], FieldInfo{
					Name:       fieldName,
					GoType:     goType,
					JSONTag:    jsonTag,
					BSONTag:    bsonTag,
					Comment:    el.Short,
					Min:        el.Min,
					MaxLength:  el.MaxLength,
					Pattern:    el.Pattern,
					Fixed:      el.Fixed,
					IsRequired: el.Min > 0,
					Path:       el.Path,
				})
			}
			continue
		}

		if isPrimitiveType && len(el.Type) > 0 {
			fhirType := el.Type[0].Code
			if def, ok := g.Definitions[fhirType]; ok && def.Kind == "primitive-type" {
				continue
			}
		}

		goType := g.mapGoType(el)

		if len(el.Type) > 0 &&
			(el.Type[0].Code == "BackboneElement" || el.Type[0].Code == "Element") {
			nestedStructName := g.deriveNestedTypeName(el.Path)
			if nestedStructName != "" {
				if isPrimitiveType {
					if def, ok := g.Definitions[nestedStructName]; ok && def.Kind == "primitive-type" {
						continue
					}
					if needsFHIRPrefix(nestedStructName) {
						continue
					}
				}
				if _, exists := structs[nestedStructName]; !exists {
					structs[nestedStructName] = []FieldInfo{}
				}
			}
		}

		cleanName := strings.ReplaceAll(lastPart, "[x]", "")
		cleanName = strings.ReplaceAll(cleanName, "-", "_")
		cleanName = text.TitleCase(cleanName)

		if cleanName == "" {
			continue
		}

		if el.Max == "*" {
			goType = "[]" + goType
		} else if el.Min == 0 && !strings.HasPrefix(goType, "[]") && goType != "json.RawMessage" && goType != "any" {
			goType = "*" + goType
		} else if el.Min > 0 && !strings.HasPrefix(goType, "[]") {
			isPrimitive := goType == "bool" || goType == "string" || goType == "int" || goType == "int64" || goType == "float64" || goType == "json.RawMessage" || goType == "any"
			if !isPrimitive {
				goType = "*" + goType
			}
		}

		escapedTag := strings.ReplaceAll(lastPart, "\"", "\\\"")
		jsonTag := fmt.Sprintf("`json:\"%s\"`", escapedTag)
		if el.Min == 0 {
			jsonTag = fmt.Sprintf("`json:\"%s,omitempty\"`", escapedTag)
		}
		bsonTag := generateBSONTag(cleanName, el.Min == 0)

		structs[structName] = append(structs[structName], FieldInfo{
			Name:       cleanName,
			GoType:     goType,
			JSONTag:    jsonTag,
			BSONTag:    bsonTag,
			Comment:    el.Short,
			Min:        el.Min,
			MaxLength:  el.MaxLength,
			Pattern:    el.Pattern,
			Fixed:      el.Fixed,
			IsRequired: el.Min > 0,
			Path:       el.Path,
		})
	}
	return structs
}

func generateBSONTag(fieldName string, hasOmitEmpty bool) string {
	snakeName := text.ToSnakeCase(fieldName)
	bsonTag := fmt.Sprintf("`bson:\"%s\"`", snakeName)
	if hasOmitEmpty {
		bsonTag = fmt.Sprintf("`bson:\"%s,omitempty\"`", snakeName)
	}
	return bsonTag
}
