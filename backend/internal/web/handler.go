package web

import (
	"html/template"
	"net/http"
	"strings"
)

type pageData struct {
	Title    string
	PageName string
}

func newPageData(r *http.Request) pageData {
	pageName := strings.Trim(r.URL.Path, "/")
	if pageName == "" {
		pageName = "home"
	}

	return pageData{
		Title:    strings.ToUpper(pageName[:1]) + pageName[1:],
		PageName: pageName,
	}
}

func handlePage(tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data := newPageData(r)

		renderTemplate(w, r, tmpl, data)
	})
}
