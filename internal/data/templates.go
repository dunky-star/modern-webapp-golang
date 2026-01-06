package data

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}

// NewTemplateData creates a new TemplateData with initialized maps
func NewTemplateData() *TemplateData {
	return &TemplateData{
		StringMap: make(map[string]string),
		IntMap:    make(map[string]int),
		FloatMap:  make(map[string]float32),
		Data:      make(map[string]interface{}),
	}
}

// SetCSRFToken sets the CSRF token (used by render package to auto-inject from context)
func (td *TemplateData) SetCSRFToken(token string) {
	td.CSRFToken = token
}
