package officegen

import (
	"encoding/json"
	"testing"
)

func TestCleanJSONResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Clean JSON",
			input:    `{"title": "Test", "slides": []}`,
			expected: `{"title": "Test", "slides": []}`,
		},
		{
			name:     "JSON with code block markers",
			input:    "```json\n{\"title\": \"Test\", \"slides\": []}\n```",
			expected: `{"title": "Test", "slides": []}`,
		},
		{
			name:     "JSON with generic code block markers",
			input:    "```\n{\"title\": \"Test\", \"slides\": []}\n```",
			expected: `{"title": "Test", "slides": []}`,
		},
		{
			name:     "JSON with leading/trailing whitespace",
			input:    "  \n  {\"title\": \"Test\", \"slides\": []}  \n  ",
			expected: `{"title": "Test", "slides": []}`,
		},
		{
			name:     "JSON with code block and extra whitespace",
			input:    "  ```json\n  {\"title\": \"Test\", \"slides\": []}  \n  ```  ",
			expected: `{"title": "Test", "slides": []}`,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only whitespace",
			input:    "   \n\n   ",
			expected: "",
		},
		{
			name:     "Nested JSON",
			input:    "```json\n{\"title\": \"Test\", \"slides\": [{\"title\": \"Slide 1\", \"bullets\": [\"Point 1\"]}]}\n```",
			expected: `{"title": "Test", "slides": [{"title": "Slide 1", "bullets": ["Point 1"]}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanJSONResponse(tt.input)
			if result != tt.expected {
				t.Errorf("cleanJSONResponse() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestPresentationJSON_Parsing(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		expectedTitle  string
		expectedSlides int
		shouldParse    bool
	}{
		{
			name:           "Valid clean JSON",
			mockResponse:   `{"title": "Test Title", "subtitle": "Test Subtitle", "slides": [{"title": "Slide 1", "bullets": ["Point 1"]}]}`,
			expectedTitle:  "Test Title",
			expectedSlides: 1,
			shouldParse:    true,
		},
		{
			name:           "JSON with markdown code blocks",
			mockResponse:   "```json\n{\"title\": \"Test Title\", \"subtitle\": \"Test Subtitle\", \"slides\": [{\"title\": \"Slide 1\", \"bullets\": [\"Point 1\"]}]}\n```",
			expectedTitle:  "Test Title",
			expectedSlides: 1,
			shouldParse:    true,
		},
		{
			name:         "Invalid JSON",
			mockResponse: "This is not JSON at all",
			shouldParse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanedResponse := cleanJSONResponse(tt.mockResponse)

			var content PresentationContent
			err := json.Unmarshal([]byte(cleanedResponse), &content)

			if tt.shouldParse {
				if err != nil {
					t.Errorf("Expected successful parse, got error: %v", err)
				}
				if content.Title != tt.expectedTitle {
					t.Errorf("Expected title %q, got %q", tt.expectedTitle, content.Title)
				}
				if len(content.Slides) != tt.expectedSlides {
					t.Errorf("Expected %d slides, got %d", tt.expectedSlides, len(content.Slides))
				}
			} else {
				if err == nil {
					t.Errorf("Expected parse error for invalid JSON, got success")
				}
			}
		})
	}
}

func TestDocumentJSON_Parsing(t *testing.T) {
	tests := []struct {
		name             string
		mockResponse     string
		expectedTitle    string
		expectedSections int
		shouldParse      bool
	}{
		{
			name:             "Valid JSON",
			mockResponse:     `{"title": "Doc Title", "sections": [{"heading": "Section 1", "content": "Content here"}]}`,
			expectedTitle:    "Doc Title",
			expectedSections: 1,
			shouldParse:      true,
		},
		{
			name:             "JSON with code blocks",
			mockResponse:     "```json\n{\"title\": \"Doc Title\", \"sections\": [{\"heading\": \"Section 1\", \"content\": \"Content here\"}]}\n```",
			expectedTitle:    "Doc Title",
			expectedSections: 1,
			shouldParse:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanedResponse := cleanJSONResponse(tt.mockResponse)

			var content DocumentContent
			err := json.Unmarshal([]byte(cleanedResponse), &content)

			if tt.shouldParse {
				if err != nil {
					t.Errorf("Expected successful parse, got error: %v", err)
				}
				if content.Title != tt.expectedTitle {
					t.Errorf("Expected title %q, got %q", tt.expectedTitle, content.Title)
				}
				if len(content.Sections) != tt.expectedSections {
					t.Errorf("Expected %d sections, got %d", tt.expectedSections, len(content.Sections))
				}
			}
		})
	}
}

func TestImagePromptsJSON_Parsing(t *testing.T) {
	tests := []struct {
		name          string
		mockResponse  string
		expectedCount int
		shouldParse   bool
	}{
		{
			name:          "Valid JSON array",
			mockResponse:  `["prompt 1", "prompt 2", "prompt 3"]`,
			expectedCount: 3,
			shouldParse:   true,
		},
		{
			name:          "JSON array with code blocks",
			mockResponse:  "```json\n[\"prompt 1\", \"prompt 2\"]\n```",
			expectedCount: 2,
			shouldParse:   true,
		},
		{
			name:         "Invalid JSON",
			mockResponse: "not an array",
			shouldParse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanedResponse := cleanJSONResponse(tt.mockResponse)

			var prompts []string
			err := json.Unmarshal([]byte(cleanedResponse), &prompts)

			if tt.shouldParse {
				if err != nil {
					t.Errorf("Expected successful parse, got error: %v", err)
				}
				if len(prompts) != tt.expectedCount {
					t.Errorf("Expected %d prompts, got %d", tt.expectedCount, len(prompts))
				}
			} else {
				if err == nil {
					t.Errorf("Expected parse error for invalid JSON, got success")
				}
			}
		})
	}
}
