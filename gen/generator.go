package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gruzdev-dev/fhir/tools/text"
)

type Generator struct {
	SpecPath    string
	OutputPath  string
	Definitions map[string]StructureDefinition
	usedTypes   map[string]bool
}

func NewGenerator(specPath, outputPath string) *Generator {
	return &Generator{
		SpecPath:    specPath,
		OutputPath:  outputPath,
		Definitions: make(map[string]StructureDefinition),
		usedTypes:   make(map[string]bool),
	}
}

func (g *Generator) Load(filename string) error {
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
		if entry.Resource.Name != "" {
			g.Definitions[entry.Resource.Name] = entry.Resource
		}
	}
	return nil
}

func (g *Generator) Generate() error {
	if err := os.MkdirAll(g.OutputPath, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(g.OutputPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read output directory: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".go" {
			if err := os.Remove(filepath.Join(g.OutputPath, entry.Name())); err != nil {
				return fmt.Errorf("remove old file %s: %w", entry.Name(), err)
			}
		}
	}

	for _, def := range g.Definitions {
		if !text.IsValidGoIdentifier(def.Name) {
			continue
		}
		if err := g.WriteResource(def); err != nil {
			return err
		}
	}
	return nil
}
