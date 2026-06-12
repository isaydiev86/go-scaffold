package generator

import (
	"fmt"
	"regexp"
	"strings"
)

// serviceNameRe allows lowercase letters, digits and single hyphens:
// starts with a letter, ends with a letter or digit (single letter is also fine).
var serviceNameRe = regexp.MustCompile(`^[a-z]([a-z0-9-]*[a-z0-9])?$`)

// ValidateServiceName checks that name is safe to use in file paths,
// Go identifiers, SQL identifiers and Kafka topic names.
func ValidateServiceName(name string) error {
	if name == "" {
		return fmt.Errorf("service name is required")
	}

	if strings.Contains(name, "--") {
		return fmt.Errorf("invalid service name %q: double hyphens are not allowed", name)
	}

	if !serviceNameRe.MatchString(name) {
		return fmt.Errorf("invalid service name %q: use lowercase letters, digits and hyphens; it must start with a letter and must not end with a hyphen (e.g. \"payments\", \"new-payments\")", name)
	}

	return nil
}

// ValidateModulePath performs a minimal sanity check of a Go module path.
func ValidateModulePath(module string) error {
	if module == "" {
		return fmt.Errorf("module path is required")
	}

	if strings.ContainsAny(module, " \t") {
		return fmt.Errorf("invalid module path %q: spaces are not allowed", module)
	}

	if strings.HasPrefix(module, "/") || strings.HasSuffix(module, "/") {
		return fmt.Errorf("invalid module path %q: must not start or end with a slash", module)
	}

	return nil
}
