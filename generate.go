package main

import (
	"fhir/gen"
	"log"
)

//go:generate go run generate.go
func main() {
	gen := gen.NewGenerator("./spec", "./r5")

	log.Println("Loading official FHIR StructureDefinitions...")

	if err := gen.Load("profiles-types.json"); err != nil {
		log.Fatal("Failed to load types:", err)
	}

	if err := gen.Load("profiles-resources.json"); err != nil {
		log.Fatal("Failed to load resources:", err)
	}

	log.Println("Generating clean Go models...")
	if err := gen.Generate(); err != nil {
		log.Fatal("Generation failed:", err)
	}

	log.Println("Done! Check 'fhir' directory.")
}
