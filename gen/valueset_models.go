package gen

import "encoding/json"

type ValueSetBundle struct {
	ResourceType string                `json:"resourceType"`
	Entry        []ValueSetBundleEntry `json:"entry"`
}

type ValueSetBundleEntry struct {
	FullUrl  string          `json:"fullUrl"`
	Resource json.RawMessage `json:"resource"`
}

type ValueSetResource struct {
	ResourceType string             `json:"resourceType"`
	ID           string             `json:"id"`
	URL          string             `json:"url"`
	Name         string             `json:"name"`
	Title        string             `json:"title,omitempty"`
	Compose      *ValueSetCompose   `json:"compose,omitempty"`
	Expansion    *ValueSetExpansion `json:"expansion,omitempty"`
}

type ValueSetCompose struct {
	Include []ValueSetComposeInclude `json:"include"`
}

type ValueSetComposeInclude struct {
	System  string                   `json:"system"`
	Concept []ValueSetComposeConcept `json:"concept,omitempty"`
}

type ValueSetComposeConcept struct {
	Code string `json:"code"`
}

type ValueSetExpansion struct {
	Contains []ValueSetExpansionContains `json:"contains"`
}

type ValueSetExpansionContains struct {
	System   string                      `json:"system,omitempty"`
	Code     string                      `json:"code"`
	Display  string                      `json:"display,omitempty"`
	Contains []ValueSetExpansionContains `json:"contains,omitempty"`
}

type CodeSystemResource struct {
	ResourceType string              `json:"resourceType"`
	ID           string              `json:"id"`
	URL          string              `json:"url"`
	Name         string              `json:"name"`
	Concept      []CodeSystemConcept `json:"concept,omitempty"`
}

type CodeSystemConcept struct {
	Code       string              `json:"code"`
	Display    string              `json:"display,omitempty"`
	Definition string              `json:"definition,omitempty"`
	Concept    []CodeSystemConcept `json:"concept,omitempty"`
}

type Code struct {
	Code       string
	Display    string
	Definition string
}
