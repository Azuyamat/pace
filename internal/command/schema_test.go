package command

import (
	"strings"
	"testing"

	"github.com/azuyamat/pace/internal/config/schema"
)

func TestLanguageReferenceIncludesSchemaEntries(t *testing.T) {
	lang := schema.Get()
	docs := languageReference(lang)

	for _, statement := range lang.TopLevel {
		if !strings.Contains(docs, "`"+statement.Name+"`") {
			t.Fatalf("generated docs missing top-level statement %q", statement.Name)
		}
	}
	for _, property := range append(append(lang.TaskProps, lang.HookProps...), lang.ArgsProps...) {
		if !strings.Contains(docs, "`"+property.Name+"`") {
			t.Fatalf("generated docs missing property %q", property.Name)
		}
	}
}

func TestGeneratedKeywordDataIncludesSchemaEntries(t *testing.T) {
	lang := schema.Get()
	keywords := generatedKeywords(lang)
	for _, statement := range lang.TopLevel {
		if !strings.Contains(keywords, statement.Name) {
			t.Fatalf("generated keyword data missing statement %q", statement.Name)
		}
	}
	for _, property := range append(append(lang.TaskProps, lang.HookProps...), lang.ArgsProps...) {
		if !strings.Contains(keywords, property.Name) {
			t.Fatalf("generated keyword data missing property %q", property.Name)
		}
	}
}
