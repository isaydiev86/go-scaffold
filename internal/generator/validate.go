package generator

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/mod/module"
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

// ValidateModulePath checks a Go module path with the same rules the go
// toolchain uses (invalid characters, reserved names, a dot in the first
// path element, etc.).
func ValidateModulePath(mod string) error {
	if mod == "" {
		return fmt.Errorf("module path is required")
	}

	// module.CheckPath errors already include the path and the reason,
	// e.g.: malformed module path "myorg/payments": missing dot in first path element
	return module.CheckPath(mod)
}
