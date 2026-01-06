package gen

import (
	"testing"
)

func TestExtractBaseType(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "pointer type",
			input: "*string",
			want:  "string",
		},
		{
			name:  "slice type",
			input: "[]string",
			want:  "string",
		},
		{
			name:  "pointer to slice",
			input: "*[]string",
			want:  "string",
		},
		{
			name:  "simple type",
			input: "string",
			want:  "string",
		},
		{
			name:  "nested pointer",
			input: "**string",
			want:  "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractBaseType(tt.input)
			if got != tt.want {
				t.Errorf("extractBaseType(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsBuiltinType(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "string is builtin",
			input: "string",
			want:  true,
		},
		{
			name:  "bool is builtin",
			input: "bool",
			want:  true,
		},
		{
			name:  "int is builtin",
			input: "int",
			want:  true,
		},
		{
			name:  "float64 is builtin",
			input: "float64",
			want:  true,
		},
		{
			name:  "any is builtin",
			input: "any",
			want:  true,
		},
		{
			name:  "json.RawMessage is builtin",
			input: "json.RawMessage",
			want:  true,
		},
		{
			name:  "custom type is not builtin",
			input: "CustomType",
			want:  false,
		},
		{
			name:  "Reference is not builtin",
			input: "Reference",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBuiltinType(tt.input)
			if got != tt.want {
				t.Errorf("isBuiltinType(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
