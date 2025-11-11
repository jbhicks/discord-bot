package officegen

type DocumentContent struct {
	Title    string            `json:"title"`
	Sections []DocumentSection `json:"sections"`
}

type DocumentSection struct {
	Heading    string   `json:"heading"`
	Content    string   `json:"content"`
	Paragraphs []string `json:"paragraphs,omitempty"`
}

type SpreadsheetContent struct {
	Title  string            `json:"title"`
	Sheets []SpreadsheetData `json:"sheets"`
}

type SpreadsheetData struct {
	Name    string     `json:"name"`
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}

type PresentationContent struct {
	Title    string              `json:"title"`
	Subtitle string              `json:"subtitle"`
	Slides   []PresentationSlide `json:"slides"`
}

type PresentationSlide struct {
	Title   string   `json:"title"`
	Bullets []string `json:"bullets"`
}

type DocumentRequest struct {
	Prompt      string
	Title       string
	TargetPages int
}

type SpreadsheetRequest struct {
	Prompt      string
	Title       string
	TargetPages int
}

type PresentationRequest struct {
	Prompt       string
	Title        string
	TargetSlides int
}

type GeneratedDocument struct {
	Data     []byte
	Filename string
}

type GeneratedSpreadsheet struct {
	Data     []byte
	Filename string
}

type GeneratedPresentation struct {
	Data     []byte
	Filename string
}
