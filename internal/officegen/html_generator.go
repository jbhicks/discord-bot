package officegen

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"strings"
)

type HTMLGenerator struct{}

func NewHTMLGenerator() *HTMLGenerator {
	return &HTMLGenerator{}
}

func (h *HTMLGenerator) GeneratePresentationHTML(content *PresentationContent, images [][]byte) (string, error) {
	slog.Info("Generating presentation HTML", "slides", len(content.Slides), "images", len(images))
	for i, slide := range content.Slides {
		slog.Info("HTML generation - slide details", "index", i, "title", slide.Title, "bullets", len(slide.Bullets))
	}

	tmpl := `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
@page {
  size: 10in 7.5in;
  margin: 0;
}
body {
  margin: 0;
  padding: 0;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}
.slide {
  width: 10in;
  height: 7.5in;
  page-break-after: always;
  position: relative;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}
.slide.title-slide {
  justify-content: center;
  align-items: center;
  text-align: center;
  padding: 2in;
}
.slide.content-slide {
  padding: 0.75in 1in;
  justify-content: space-between;
}
.content-container {
  flex: 0 0 auto;
}
.title-slide h1 {
  font-size: 60px;
  margin: 0 0 30px 0;
  font-weight: 700;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}
.title-slide h2 {
  font-size: 32px;
  margin: 0;
  font-weight: 300;
  opacity: 0.9;
}
.content-slide h2 {
  font-size: 44px;
  margin: 0 0 30px 0;
  font-weight: 600;
  border-bottom: 3px solid rgba(255,255,255,0.3);
  padding-bottom: 15px;
}
.content-slide ul {
  list-style: none;
  padding: 0;
  margin: 0;
  font-size: 24px;
  line-height: 1.6;
}
.content-slide li {
  margin: 15px 0;
  padding-left: 40px;
  position: relative;
}
.content-slide li:before {
  content: "▸";
  position: absolute;
  left: 0;
  color: rgba(255,255,255,0.8);
  font-size: 28px;
}
.slide-image {
  flex: 0 0 auto;
  text-align: center;
  max-height: 3.5in;
  overflow: hidden;
}
.slide-image img {
  max-width: 100%;
  max-height: 3.5in;
  height: auto;
  width: auto;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
}
</style>
</head>
<body>
<div class="slide title-slide">
  <h1>{{.Title}}</h1>
  <h2>{{.Subtitle}}</h2>
</div>
{{range $idx, $slide := .Slides}}
<div class="slide content-slide">
  <div class="content-container">
    <h2>{{if $slide.Title}}{{$slide.Title}}{{else}}Slide {{add $idx 1}}{{end}}</h2>
    {{if $slide.Bullets}}
    <ul>
      {{range $slide.Bullets}}
      <li>{{.}}</li>
      {{end}}
    </ul>
    {{else}}
    <p style="font-size: 20px; opacity: 0.8; margin-top: 50px;">No content available for this slide.</p>
    {{end}}
  </div>
  {{if and (lt $idx (len $.Images)) (index $.Images $idx)}}
  <div class="slide-image">
    <img src="data:image/png;base64,{{index $.Images $idx}}" />
  </div>
  {{end}}
</div>
{{end}}
</body>
</html>`

	data := struct {
		Title    string
		Subtitle string
		Slides   []PresentationSlide
		Images   []string
	}{
		Title:    content.Title,
		Subtitle: content.Subtitle,
		Slides:   content.Slides,
		Images:   make([]string, len(images)),
	}

	for i, img := range images {
		encoded := base64.StdEncoding.EncodeToString(img)
		data.Images[i] = encoded
		slog.Info("Encoded image for template", "index", i, "size", len(img), "base64_length", len(encoded))
	}

	slog.Info("Template data prepared", "title", data.Title, "subtitle", data.Subtitle, "slides", len(data.Slides), "images", len(data.Images))

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}

	t, err := template.New("presentation").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	htmlOutput := buf.String()
	slog.Info("Presentation HTML generated", "length", len(htmlOutput))

	// Log a snippet to verify content
	if len(htmlOutput) > 500 {
		slog.Info("HTML snippet", "first_500_chars", htmlOutput[:500])
	}

	// Save HTML to file for debugging
	if err := os.WriteFile("/tmp/presentation_debug.html", []byte(htmlOutput), 0644); err != nil {
		slog.Warn("Failed to write debug HTML file", "error", err)
	} else {
		slog.Info("Debug HTML written to /tmp/presentation_debug.html")
	}

	return htmlOutput, nil
}

func (h *HTMLGenerator) GenerateDocumentHTML(content *DocumentContent, images [][]byte) (string, error) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
@page {
  size: A4;
  margin: 1in;
}
body {
  font-family: 'Georgia', serif;
  font-size: 12pt;
  line-height: 1.6;
  color: #333;
  margin: 0;
  padding: 0;
}
h1 {
  font-size: 28pt;
  margin: 0 0 20pt 0;
  color: #2c3e50;
  text-align: center;
  border-bottom: 2pt solid #3498db;
  padding-bottom: 10pt;
}
h2 {
  font-size: 20pt;
  margin: 30pt 0 15pt 0;
  color: #34495e;
  border-bottom: 1pt solid #bdc3c7;
  padding-bottom: 5pt;
}
h3 {
  font-size: 16pt;
  margin: 20pt 0 10pt 0;
  color: #34495e;
}
p {
  margin: 10pt 0;
  text-align: justify;
}
.image-container {
  text-align: center;
  margin: 20pt 0;
}
.image-container img {
  max-width: 100%;
  max-height: 4in;
  border: 1pt solid #bdc3c7;
  border-radius: 4px;
}
</style>
</head>
<body>
<h1>{{.Title}}</h1>
{{range $idx, $section := .Sections}}
<h2>{{$section.Heading}}</h2>
{{if $section.Paragraphs}}
  {{range $section.Paragraphs}}
<p>{{.}}</p>
  {{end}}
{{else if $section.Content}}
<p>{{$section.Content}}</p>
{{end}}
{{if and (lt $idx (len $.Images)) (index $.Images $idx)}}
<div class="image-container">
  <img src="data:image/png;base64,{{index $.Images $idx}}" />
</div>
{{end}}
{{end}}
</body>
</html>`

	data := struct {
		Title    string
		Sections []DocumentSection
		Images   []string
	}{
		Title:    content.Title,
		Sections: content.Sections,
		Images:   make([]string, len(images)),
	}

	for i, img := range images {
		data.Images[i] = base64.StdEncoding.EncodeToString(img)
	}

	t, err := template.New("document").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (h *HTMLGenerator) GenerateSpreadsheetHTML(content *SpreadsheetContent) (string, error) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
@page {
  size: A4 landscape;
  margin: 0.5in;
}
body {
  font-family: 'Arial', sans-serif;
  font-size: 10pt;
  margin: 0;
  padding: 0;
}
h1 {
  font-size: 20pt;
  margin: 0 0 20pt 0;
  color: #2c3e50;
  text-align: center;
}
h2 {
  font-size: 14pt;
  margin: 30pt 0 10pt 0;
  color: #34495e;
}
table {
  width: 100%;
  border-collapse: collapse;
  margin: 10pt 0 30pt 0;
  font-size: 9pt;
}
thead {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}
th {
  padding: 8pt;
  text-align: left;
  font-weight: 600;
  border: 1pt solid #fff;
}
td {
  padding: 6pt 8pt;
  border: 1pt solid #ddd;
}
tbody tr:nth-child(even) {
  background: #f8f9fa;
}
tbody tr:hover {
  background: #e9ecef;
}
.number {
  text-align: right;
}
</style>
</head>
<body>
<h1>{{.Title}}</h1>
{{range .Sheets}}
<h2>{{.Name}}</h2>
<table>
<thead>
<tr>
{{range .Headers}}
<th>{{.}}</th>
{{end}}
</tr>
</thead>
<tbody>
{{range .Rows}}
<tr>
{{range .}}
<td class="{{if isNumber .}}number{{end}}">{{.}}</td>
{{end}}
</tr>
{{end}}
</tbody>
</table>
{{end}}
</body>
</html>`

	funcMap := template.FuncMap{
		"isNumber": func(s string) bool {
			s = strings.TrimSpace(s)
			if s == "" {
				return false
			}
			for _, c := range s {
				if (c < '0' || c > '9') && c != '.' && c != '-' && c != ',' && c != '$' && c != '%' {
					return false
				}
			}
			return true
		},
	}

	t, err := template.New("spreadsheet").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, content); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	htmlOutput := buf.String()
	slog.Info("Spreadsheet HTML generated", "length", len(htmlOutput))

	// Log a snippet to verify content
	if len(htmlOutput) > 500 {
		slog.Info("HTML snippet", "first_500_chars", htmlOutput[:500])
	}

	return htmlOutput, nil
}
