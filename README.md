# FHIR R5 Generator

A Go code generator for FHIR R5 resources and types. Generates type-safe Go models from official FHIR StructureDefinitions.

## Overview

This generator reads FHIR R5 specification files and generates complete Go struct definitions for all resources and complex types defined in the specification. Each generated model includes JSON/BSON tags and validation methods.

## Generated Models

Generated models are located in the `r5/` directory. The generator creates one Go file per resource/type with:

- Struct definitions matching FHIR specification
- JSON and BSON tags for serialization
- `Validate()` methods for field validation
- Proper handling of required fields, cardinality, patterns, and constraints

## Usage

### Generating Models

Run the generator to create all FHIR R5 models:

```bash
go run generate.go
```

This will:
1. Load StructureDefinitions from `spec/profiles-types.json` and `spec/profiles-resources.json`
2. Generate Go models for all resources and types
3. Write output files to `r5/` directory

### Using Generated Models

Import the generated models in your code:

```go
import "fhir/r5"

// Use generated types
patient := &r5.Patient{
    ID: stringPtr("example"),
    // ... other fields
}

// Validate the resource
if err := patient.Validate(); err != nil {
    // handle validation error
}

// Marshal to JSON
data, err := json.Marshal(patient)
```

## Requirements

- Go 1.25 or later
- FHIR R5 specification files in `spec/` directory
