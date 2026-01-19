package data

import "github.com/dunky-star/modern-webapp-golang/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated bool
}

// SetCSRFToken sets the CSRF token (used by render package to auto-inject from context)
func (td *TemplateData) SetCSRFToken(token string) {
	td.CSRFToken = token
}
