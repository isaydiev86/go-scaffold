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
	skip := skippedPaths(vars)

	return fs.WalkDir(templates.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root and the embed.go file itself
		if path == "." || path == "embed.go" {
			return nil
		}

		// Skip templates of integrations that were not selected
		if isSkipped(path, skip) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		// Build output path: replace the "__service__" placeholder → ServiceName, strip .tmpl
		outRel := strings.ReplaceAll(path, "__service__", vars.ServiceName)
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

// skippedPaths returns embedded template paths (files or whole directories)
// that must not be generated for the given options.
func skippedPaths(vars Vars) []string {
	var skip []string

	if !vars.WithRedis {
		skip = append(skip,
			"pkg/redis",
			"internal/adapter/redis",
		)
	}

	if !vars.WithKafkaConsumer {
		skip = append(skip, "internal/controller/kafka_consumer")
	}

	if !vars.WithKafkaProducer {
		skip = append(skip,
			"internal/adapter/kafka_producer",
			"internal/controller/worker",
			"internal/domain/event.go.tmpl",
			"internal/usecase/outbox_read_and_produce.go.tmpl",
			"internal/adapter/postgres/save_outbox_kafka.go.tmpl",
			"internal/adapter/postgres/read_outbox_kafka.go.tmpl",
			"migration/postgres/20240101000002_outbox.up.sql.tmpl",
			"migration/postgres/20240101000002_outbox.down.sql.tmpl",
		)
	}

	if !vars.WithKafka {
		skip = append(skip, "pkg/logger/kafka.go.tmpl")
	}

	return skip
}

// isSkipped reports whether path matches one of the skipped paths exactly
// or lies inside a skipped directory.
func isSkipped(path string, skip []string) bool {
	for _, s := range skip {
		if path == s || strings.HasPrefix(path, s+"/") {
			return true
		}
	}

	return false
}
