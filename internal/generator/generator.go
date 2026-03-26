package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	templates "github.com/isaydiev86/go-scaffold/templates"
)

// Generate renders all templates into outputDir, substituting Vars.
func Generate(vars Vars, outputDir string) error {
	return fs.WalkDir(templates.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root and the embed.go file itself
		if path == "." || path == "embed.go" {
			return nil
		}

		// Build output path: replace "service" placeholder → ServiceName, strip .tmpl
		outRel := strings.ReplaceAll(path, "service", vars.ServiceName)
		outRel = strings.TrimSuffix(outRel, ".tmpl")
		outPath := filepath.Join(outputDir, outRel)

		if d.IsDir() {
			return os.MkdirAll(outPath, 0o755)
		}

		// Only process .tmpl files
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		content, err := templates.FS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read template %s: %w", path, err)
		}

		tmpl, err := template.New(path).Parse(string(content))
		if err != nil {
			return fmt.Errorf("parse template %s: %w", path, err)
		}

		if err = os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(outPath), err)
		}

		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create %s: %w", outPath, err)
		}
		defer f.Close()

		if err = tmpl.Execute(f, vars); err != nil {
			return fmt.Errorf("execute template %s: %w", path, err)
		}

		return nil
	})
}