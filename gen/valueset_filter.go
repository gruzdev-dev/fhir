package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (g *Generator) ExtractRequiredValueSets() (map[string]bool, error) {
	requiredURLs := make(map[string]bool)

	if err := g.extractRequiredFromFile("profiles-resources.json", requiredURLs); err != nil {
		return nil, fmt.Errorf("extract from profiles-resources.json: %w", err)
	}

	if err := g.extractRequiredFromFile("profiles-types.json", requiredURLs); err != nil {
		return nil, fmt.Errorf("extract from profiles-types.json: %w", err)
	}

	return requiredURLs, nil
}

func (g *Generator) extractRequiredFromFile(filename string, requiredURLs map[string]bool) error {
	path := filepath.Join(g.SpecPath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var bundle StructureDefinitionBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return fmt.Errorf("unmarshal %s: %w", filename, err)
	}

	for _, entry := range bundle.Entry {
		for _, element := range entry.Resource.Snapshot.Element {
			if element.Binding != nil && element.Binding.Strength == "required" {
				if element.Binding.ValueSet != "" {
					url := normalizeValueSetURL(element.Binding.ValueSet)
					requiredURLs[url] = true
				}
			}
		}
	}

	return nil
}

func normalizeValueSetURL(url string) string {
	parts := strings.Split(url, "|")
	return parts[0]
}
