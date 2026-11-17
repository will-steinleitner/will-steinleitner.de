package web

import (
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data any) {
	templateName := "base"

	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
