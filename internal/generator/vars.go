package generator

import (
	"strings"
	"unicode"
)

// Options toggles optional integrations in the generated service.
type Options struct {
	WithRedis         bool // cache adapter + pkg/redis + docker-compose service
	WithKafkaConsumer bool // kafka_consumer controller
	WithKafkaProducer bool // kafka_producer adapter + transactional outbox + worker
}

// Vars holds all template variables for service generation.
type Vars struct {
	ServiceName       string // e.g. "new-payments" — lowercase, hyphens preserved (for paths/URLs)
	ServiceNameSQL    string // e.g. "new_payments" — hyphens replaced with underscores (for SQL identifiers)
	ServiceNameTitle  string // e.g. "New-Payments" — each word title-cased
	EntityName        string // e.g. "NewPayment" — singular PascalCase for Go structs/methods
	ModuleName        string // e.g. "github.com/myorg/new-payments"
	GoVersion         string // e.g. "1.26"
	WithRedis         bool   // include Redis (cache adapter, pkg/redis, docker-compose service)
	WithKafkaConsumer bool   // include Kafka consumer controller
	WithKafkaProducer bool   // include Kafka producer adapter + outbox worker
	WithKafka         bool   // any Kafka part is enabled (broker in docker-compose, kafka-go dependency)
}

// NewVars builds Vars from a service name and module path.
func NewVars(serviceName, moduleName string, opts Options) Vars {
	name := strings.ToLower(serviceName)
	return Vars{
		ServiceName:       name,
		ServiceNameSQL:    strings.ReplaceAll(name, "-", "_"),
		ServiceNameTitle:  titleWords(name),
		EntityName:        singular(toPascalCase(name)),
		ModuleName:        moduleName,
		GoVersion:         "1.26",
		WithRedis:         opts.WithRedis,
		WithKafkaConsumer: opts.WithKafkaConsumer,
		WithKafkaProducer: opts.WithKafkaProducer,
		WithKafka:         opts.WithKafkaConsumer || opts.WithKafkaProducer,
	}
}

// toPascalCase converts a hyphen/underscore-separated string to PascalCase.
// e.g. "new-payment" → "NewPayment", "payments" → "Payments"
func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_'
	})
	var b strings.Builder
	for _, p := range parts {
		if len(p) > 0 {
			r := []rune(p)
			r[0] = unicode.ToUpper(r[0])
			b.WriteString(string(r))
		}
	}
	return b.String()
}

// titleWords capitalises the first letter of each hyphen-separated word.
// e.g. "new-payments" → "New-Payments"
func titleWords(s string) string {
	parts := strings.Split(s, "-")
	for i, p := range parts {
		if len(p) > 0 {
			r := []rune(p)
			r[0] = unicode.ToUpper(r[0])
			parts[i] = string(r)
		}
	}
	return strings.Join(parts, "-")
}

// singular removes trailing 's' to produce a singular form.
func singular(s string) string {
	if strings.HasSuffix(s, "s") && len(s) > 1 {
		return s[:len(s)-1]
	}

	return s
}
