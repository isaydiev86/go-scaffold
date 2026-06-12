//go:build smoke

package generator

import (
	"os/exec"
	"path/filepath"
	"testing"
)

// TestSmokeBootstrap generates a project and runs the documented bootstrap
// sequence (make bootstrap → go test ./...). It needs network access and
// mockery + oapi-codegen in PATH, hence the "smoke" build tag:
//
//	go test -tags=smoke -timeout 20m ./internal/generator
func TestSmokeBootstrap(t *testing.T) {
	for _, tool := range []string{"mockery", "oapi-codegen", "make"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not installed: %v", tool, err)
		}
	}

	cases := []struct {
		name string
		opts Options
	}{
		{name: "base", opts: Options{}},
		{name: "all", opts: Options{WithRedis: true, WithKafkaConsumer: true, WithKafkaProducer: true}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outDir := filepath.Join(t.TempDir(), "payments")

			vars := NewVars("payments", "github.com/myorg/payments", tc.opts)
			if err := Generate(vars, outDir); err != nil {
				t.Fatalf("Generate: %v", err)
			}

			run(t, outDir, "make", "bootstrap")
			run(t, outDir, "go", "test", "./...")
		})
	}
}

func run(t *testing.T, dir string, name string, args ...string) {
	t.Helper()

	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
	}
}
