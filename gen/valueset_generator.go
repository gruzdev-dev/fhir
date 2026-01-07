package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (g *Generator) LoadValueSets() (map[string]ValueSetResource, map[string]CodeSystemResource, error) {
	path := filepath.Join(g.SpecPath, "valuesets.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read valuesets.json: %w", err)
	}

	var bundle ValueSetBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, nil, fmt.Errorf("unmarshal valuesets.json: %w", err)
	}

	valueSets := make(map[string]ValueSetResource)
	codeSystems := make(map[string]CodeSystemResource)

	for _, entry := range bundle.Entry {
		var resourceType struct {
			ResourceType string `json:"resourceType"`
		}
		if err := json.Unmarshal(entry.Resource, &resourceType); err != nil {
			continue
		}

		switch resourceType.ResourceType {
		case "ValueSet":
			var vs ValueSetResource
			if err := json.Unmarshal(entry.Resource, &vs); err != nil {
				continue
			}
			if vs.URL != "" {
				valueSets[vs.URL] = vs
			}

		case "CodeSystem":
			var cs CodeSystemResource
			if err := json.Unmarshal(entry.Resource, &cs); err != nil {
				continue
			}
			if cs.URL != "" {
				codeSystems[cs.URL] = cs
			}
		}
	}

	return valueSets, codeSystems, nil
}
