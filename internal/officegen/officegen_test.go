package officegen

import (
	"strings"
	"testing"
)

func TestDocumentGenerator_BasicGeneration(t *testing.T) {
	t.Skip("Integration test - requires LLM server")

	client := NewClient("http://localhost:8081")
	gen := NewDocumentGenerator(client, nil)

	req := &DocumentRequest{
		Prompt:      "Write about artificial intelligence",
		Title:       "AI Overview",
		TargetPages: 1,
	}

	result, err := gen.Generate(req)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	if !strings.HasSuffix(result.Filename, ".pdf") {
		t.Errorf("Expected .pdf extension, got: %s", result.Filename)
	}

	if len(result.Data) < 1000 {
		t.Errorf("Document seems too small: %d bytes", len(result.Data))
	}
}

func TestSpreadsheetGenerator_BasicGeneration(t *testing.T) {
	t.Skip("Integration test - requires LLM server")

	client := NewClient("http://localhost:8081")
	gen := NewSpreadsheetGenerator(client, nil)

	req := &SpreadsheetRequest{
		Prompt:      "Create a budget spreadsheet",
		Title:       "Monthly Budget",
		TargetPages: 1,
	}

	result, err := gen.Generate(req)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	if !strings.HasSuffix(result.Filename, ".pdf") {
		t.Errorf("Expected .pdf extension, got: %s", result.Filename)
	}

	if len(result.Data) < 1000 {
		t.Errorf("Spreadsheet seems too small: %d bytes", len(result.Data))
	}
}

func TestPresentationGenerator_BasicGeneration(t *testing.T) {
	t.Skip("Integration test - requires LLM server")

	client := NewClient("http://localhost:8081")
	gen := NewPresentationGenerator(client, nil)

	req := &PresentationRequest{
		Prompt:       "Create a presentation about solar energy",
		Title:        "Solar Energy",
		TargetSlides: 3,
	}

	result, err := gen.Generate(req)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	if !strings.HasSuffix(result.Filename, ".pdf") {
		t.Errorf("Expected .pdf extension, got: %s", result.Filename)
	}

	if len(result.Data) < 10000 {
		t.Errorf("Presentation seems too small: %d bytes", len(result.Data))
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple title", "Simple Title", "Simple_Title"},
		{"Special chars", "Title with special !@#$ chars", "Title_with_special_chars"},
		{"Dashes", "Title-with-dashes", "Title-with-dashes"},
		{"Underscores", "Title_with_underscores", "Title_with_underscores"},
		{"Extra spaces", "   Spaces   ", "Spaces"},
		{"Empty", "", "document"},
		{"Too long", strings.Repeat("a", 150), strings.Repeat("a", 100)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPDFExporter_CheckAvailable(t *testing.T) {
	exporter := NewPDFExporter()
	available := exporter.CheckAvailable()

	t.Logf("LibreOffice available: %v", available)

	if !available {
		t.Skip("LibreOffice not available, skipping PDF tests")
	}
}
