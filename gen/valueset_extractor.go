package gen

func ExtractCodesFromValueSet(vs ValueSetResource, codeSystems map[string]CodeSystemResource) ([]Code, error) {
	var codes []Code

	if vs.Expansion != nil && len(vs.Expansion.Contains) > 0 {
		codes = flattenExpansion(vs.Expansion.Contains)
		return codes, nil
	}

	if vs.Compose != nil && len(vs.Compose.Include) > 0 {
		for _, include := range vs.Compose.Include {
			if cs, exists := codeSystems[include.System]; exists {
				if len(cs.Concept) > 0 {
					codes = append(codes, flattenConcepts(cs.Concept)...)
				}
			}
			if len(include.Concept) > 0 {
				for _, concept := range include.Concept {
					codes = append(codes, Code{
						Code: concept.Code,
					})
				}
			}
		}
	}

	return codes, nil
}

func flattenConcepts(concepts []CodeSystemConcept) []Code {
	var codes []Code
	for _, concept := range concepts {
		codes = append(codes, Code{
			Code:       concept.Code,
			Display:    concept.Display,
			Definition: concept.Definition,
		})
		if len(concept.Concept) > 0 {
			codes = append(codes, flattenConcepts(concept.Concept)...)
		}
	}
	return codes
}

func flattenExpansion(contains []ValueSetExpansionContains) []Code {
	var codes []Code
	for _, item := range contains {
		codes = append(codes, Code{
			Code:       item.Code,
			Display:    item.Display,
			Definition: "",
		})
		if len(item.Contains) > 0 {
			codes = append(codes, flattenExpansion(item.Contains)...)
		}
	}
	return codes
}
