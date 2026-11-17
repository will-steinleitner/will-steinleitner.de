package web

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"slices"

	"github.com/will-steinleitner/will-steinleitner.de/internal/common/server/middleware"
)

//go:embed assets/* all:templates
var content embed.FS

func AddRoute(
	mux *http.ServeMux,
	ctx context.Context,
	logger *slog.Logger,
) error {
	templateFiles, err := templateFiles(content)
	if err != nil {
		return err
	}

	tmpl := template.Must(template.ParseFS(content, templateFiles...))

	assets, err := fs.Sub(content, "assets")
	if err != nil {
		return fmt.Errorf("load web assets: %w", err)
	}

	defaultMiddle := middleware.NewDefault(logger)

	pageTemplate := func(name string) *template.Template {
		pageTemplate, err := pageTemplateSet(tmpl, name)
		if err != nil {
			panic(fmt.Errorf("load %s page template: %w", name, err))
		}

		return pageTemplate
	}

	// Pages
	homeTmpl := pageTemplate("home/page")
	mux.Handle(http.MethodGet+" /{$}", defaultMiddle(handlePage(homeTmpl)))

	// Static content
	mux.Handle(http.MethodGet+" /assets/", defaultMiddle(http.StripPrefix("/assets/", http.FileServerFS(assets))))

	return nil
}

func templateFiles(content fs.FS) ([]string, error) {
	files := make([]string, 0)
	err := fs.WalkDir(content, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case ".html", ".gohtml":
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("list web templates: %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("list web templates: no supported templates found")
	}

	slices.Sort(files)

	return files, nil
}

func pageTemplateSet(base *template.Template, pageTemplate string) (*template.Template, error) {
	cloned, err := base.Clone()
	if err != nil {
		return nil, fmt.Errorf("clone web templates: %w", err)
	}

	if _, err := cloned.Parse(`{{ define "page-body" }}{{ template "` + pageTemplate + `" . }}{{ end }}`); err != nil {
		return nil, fmt.Errorf("parse page template alias: %w", err)
	}

	return cloned, nil
}
