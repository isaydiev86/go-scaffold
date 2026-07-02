package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// optional template outputs per feature (paths inside the generated project,
// already with "service" → "payments" and ".tmpl" stripped).
var (
	redisFiles = []string{
		"pkg/redis/redis.go",
		"internal/adapter/redis/redis.go",
		"internal/adapter/redis/get.go",
		"internal/adapter/redis/set.go",
		"internal/adapter/redis/delete.go",
	}
	consumerFiles = []string{
		"internal/controller/kafka_consumer/consumer.go",
	}
	producerFiles = []string{
		"internal/adapter/kafka_producer/producer.go",
		"internal/adapter/kafka_producer/produce.go",
		"internal/controller/worker/outbox_kafka.go",
		"internal/domain/event.go",
		"internal/usecase/outbox_read_and_produce.go",
		"internal/adapter/postgres/save_outbox_kafka.go",
		"internal/adapter/postgres/read_outbox_kafka.go",
		"migration/postgres/20240101000002_outbox.up.sql",
		"migration/postgres/20240101000002_outbox.down.sql",
	}
	kafkaFiles = []string{
		"pkg/logger/kafka.go",
	}
)

func TestGenerate(t *testing.T) {
	cases := []struct {
		name string
		opts Options
	}{
		{name: "base", opts: Options{}},
		{name: "redis", opts: Options{WithRedis: true}},
		{name: "kafka-consumer", opts: Options{WithKafkaConsumer: true}},
		{name: "kafka-producer", opts: Options{WithKafkaProducer: true}},
		{name: "all", opts: Options{WithRedis: true, WithKafkaConsumer: true, WithKafkaProducer: true}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outDir := filepath.Join(t.TempDir(), "payments")

			vars := NewVars("payments", "github.com/myorg/payments", tc.opts)
			if err := Generate(vars, outDir); err != nil {
				t.Fatalf("Generate: %v", err)
			}

			// Base files must always be present.
			assertExists(t, outDir, []string{
				"go.mod",
				"Makefile",
				"cmd/app/main.go",
				"internal/usecase/usecase.go",
				"api/http/payments_v1.yaml",
				"gen/http/payments_v1/server/generate.go",
				"pkg/payments_client_gen/client.go",
			})

			assertOptional(t, outDir, redisFiles, tc.opts.WithRedis)
			assertOptional(t, outDir, consumerFiles, tc.opts.WithKafkaConsumer)
			assertOptional(t, outDir, producerFiles, tc.opts.WithKafkaProducer)
			assertOptional(t, outDir, kafkaFiles, tc.opts.WithKafkaConsumer || tc.opts.WithKafkaProducer)

			// No generated path may keep the "__service__" placeholder or the .tmpl suffix.
			err := filepath.WalkDir(outDir, func(path string, _ os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				rel, _ := filepath.Rel(outDir, path)
				if strings.Contains(rel, "__service__") {
					t.Errorf("path %q still contains the \"__service__\" placeholder", rel)
				}
				if strings.HasSuffix(rel, ".tmpl") {
					t.Errorf("path %q still has the .tmpl suffix", rel)
				}
				return nil
			})
			if err != nil {
				t.Fatalf("walk output: %v", err)
			}
		})
	}
}

// assertOptional checks that all files exist when enabled and none exist when disabled.
func assertOptional(t *testing.T, outDir string, files []string, enabled bool) {
	t.Helper()

	if enabled {
		assertExists(t, outDir, files)
		return
	}

	for _, f := range files {
		if _, err := os.Stat(filepath.Join(outDir, f)); err == nil {
			t.Errorf("file %q must not be generated for these options", f)
		}
	}
}

func assertExists(t *testing.T, outDir string, files []string) {
	t.Helper()

	for _, f := range files {
		if _, err := os.Stat(filepath.Join(outDir, f)); err != nil {
			t.Errorf("expected file %q: %v", f, err)
		}
	}
}
