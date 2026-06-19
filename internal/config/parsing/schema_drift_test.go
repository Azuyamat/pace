package parsing

import (
	"testing"

	"github.com/azuyamat/pace/internal/config/schema"
)

func TestSchemaMatchesParserRegistries(t *testing.T) {
	lang := schema.Get()

	for _, statement := range lang.TopLevel {
		if _, ok := statementRegistry[statement.Name]; !ok {
			t.Fatalf("schema top-level statement %q is missing from parser registry", statement.Name)
		}
	}
	if len(statementRegistry) != len(lang.TopLevel) {
		t.Fatalf("parser top-level registry has %d entries, schema has %d", len(statementRegistry), len(lang.TopLevel))
	}

	for _, property := range lang.TaskProps {
		if _, ok := taskPropertyRegistry[property.Name]; !ok {
			t.Fatalf("schema task property %q is missing from parser registry", property.Name)
		}
	}
	if len(taskPropertyRegistry) != len(lang.TaskProps) {
		t.Fatalf("parser task property registry has %d entries, schema has %d", len(taskPropertyRegistry), len(lang.TaskProps))
	}

	for _, property := range lang.HookProps {
		if _, ok := hookPropertyRegistry[property.Name]; !ok {
			t.Fatalf("schema hook property %q is missing from parser registry", property.Name)
		}
	}
	if len(hookPropertyRegistry) != len(lang.HookProps) {
		t.Fatalf("parser hook property registry has %d entries, schema has %d", len(hookPropertyRegistry), len(lang.HookProps))
	}
}

func TestGeneratedSnippetsHaveContent(t *testing.T) {
	lang := schema.Get()
	for _, statement := range lang.TopLevel {
		if statement.Snippet == "" || statement.Description == "" {
			t.Fatalf("top-level statement %q must have snippet and docs", statement.Name)
		}
	}
	for _, property := range append(append(lang.TaskProps, lang.HookProps...), lang.ArgsProps...) {
		if property.Snippet == "" || property.Description == "" || property.Type == "" {
			t.Fatalf("property %q must have type, snippet, and docs", property.Name)
		}
	}
}
