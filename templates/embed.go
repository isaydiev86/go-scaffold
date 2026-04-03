package templates

import "embed"

//go:embed all:api all:cmd all:config all:gen all:internal all:migration all:pkg all:test Makefile.tmpl Dockerfile.tmpl docker-compose.yml.tmpl go.mod.tmpl .env.example.tmpl .gitignore.tmpl .golangci.yml.tmpl .mockery.yml.tmpl
var FS embed.FS