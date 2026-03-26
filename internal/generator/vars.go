package generator

import (
	"strings"
	"unicode"
)

// Vars holds all template variables for service generation.
type Vars struct {
	ServiceName      string // e.g. "payments" — lowercase
	ServiceNameTitle string // e.g. "Payments" — title case
	EntityName       string // e.g. "Payment" — singular PascalCase for structs
	ModuleName       string // e.g. "github.com/myorg/payments"
	GoVersion        string // e.g. "1.26"
}

// NewVars builds Vars from a service name and module path.
func NewVars(serviceName, moduleName string) Vars {
	name := strings.ToLower(serviceName)
	return Vars{
		ServiceName:      name,
		ServiceNameTitle: title(name),
		EntityName:       singular(title(name)),
		ModuleName:       moduleName,
		GoVersion:        "1.26",
	}
}

func title(s string) string {
	if s == "" {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return string(r)
}

// singular removes trailing 's' to produce a singular form.
func singular(s string) string {
	if strings.HasSuffix(s, "s") && len(s) > 1 {
		return s[:len(s)-1]
	}

	return s
}