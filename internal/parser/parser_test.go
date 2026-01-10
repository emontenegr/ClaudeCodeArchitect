package parser

import "testing"

func TestFindSection(t *testing.T) {
	structure := &SpecStructure{
		Sections: []SectionInfo{
			{Title: "API Endpoints", FilePath: "api.adoc", Level: 1},
			{Title: "User Types", FilePath: "types.adoc", Level: 1},
			{Title: "Database Schema", FilePath: "db/schema.adoc", Level: 1},
			{Title: "Performance", FilePath: "perf.adoc", Level: 1},
		},
	}

	tests := []struct {
		name        string
		query       string
		expectTitle string
		expectNil   bool
	}{
		{
			name:        "exact match",
			query:       "API Endpoints",
			expectTitle: "API Endpoints",
		},
		{
			name:        "case insensitive exact match",
			query:       "api endpoints",
			expectTitle: "API Endpoints",
		},
		{
			name:        "partial match",
			query:       "User",
			expectTitle: "User Types",
		},
		{
			name:        "file path match",
			query:       "schema.adoc",
			expectTitle: "Database Schema",
		},
		{
			name:      "no match",
			query:     "nonexistent",
			expectNil: true,
		},
		{
			name:        "whitespace handling",
			query:       "  Performance  ",
			expectTitle: "Performance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindSection(structure, tt.query)

			if tt.expectNil {
				if result != nil {
					t.Errorf("expected nil, got section %q", result.Title)
				}
				return
			}

			if result == nil {
				t.Errorf("expected section %q, got nil", tt.expectTitle)
				return
			}

			if result.Title != tt.expectTitle {
				t.Errorf("expected %q, got %q", tt.expectTitle, result.Title)
			}
		})
	}
}

func TestFindSectionsByFile(t *testing.T) {
	structure := &SpecStructure{
		Sections: []SectionInfo{
			{Title: "Section 1", FilePath: "a.adoc"},
			{Title: "Section 2", FilePath: "a.adoc"},
			{Title: "Section 3", FilePath: "b.adoc"},
			{Title: "Section 4", FilePath: "a.adoc"},
		},
	}

	result := FindSectionsByFile(structure, "a.adoc")

	if len(result) != 3 {
		t.Errorf("expected 3 sections, got %d", len(result))
	}

	result = FindSectionsByFile(structure, "nonexistent.adoc")
	if len(result) != 0 {
		t.Errorf("expected 0 sections, got %d", len(result))
	}
}

func TestSectionPattern(t *testing.T) {
	tests := []struct {
		line        string
		expectMatch bool
		expectLevel int
		expectTitle string
	}{
		{"= Document Title", true, 0, "Document Title"},
		{"== Section Level 1", true, 1, "Section Level 1"},
		{"=== Subsection", true, 2, "Subsection"},
		{"==== Deep Section", true, 3, "Deep Section"},
		{"Not a section", false, 0, ""},
		{"=No space", false, 0, ""},
		{"  == Indented", false, 0, ""},
		{"== Title with = sign", true, 1, "Title with = sign"},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			matches := sectionPattern.FindStringSubmatch(tt.line)

			if tt.expectMatch {
				if matches == nil {
					t.Errorf("expected match, got nil")
					return
				}
				level := len(matches[1]) - 1
				title := matches[2]

				if level != tt.expectLevel {
					t.Errorf("expected level %d, got %d", tt.expectLevel, level)
				}
				if title != tt.expectTitle {
					t.Errorf("expected title %q, got %q", tt.expectTitle, title)
				}
			} else {
				if matches != nil {
					t.Errorf("expected no match, got %v", matches)
				}
			}
		})
	}
}
