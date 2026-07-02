package generator

import "testing"

func TestValidateServiceName(t *testing.T) {
	valid := []string{
		"payments",
		"new-payments",
		"a",
		"svc2",
		"a1-b2-c3",
	}
	for _, name := range valid {
		if err := ValidateServiceName(name); err != nil {
			t.Errorf("ValidateServiceName(%q) = %v, want nil", name, err)
		}
	}

	invalid := []string{
		"",
		"Payments",
		"new payments",
		"new_payments",
		"new/payments",
		"payments.api",
		"new--payments",
		"payments-",
		"-payments",
		"1payments",
		"платежи",
	}
	for _, name := range invalid {
		if err := ValidateServiceName(name); err == nil {
			t.Errorf("ValidateServiceName(%q) = nil, want error", name)
		}
	}
}

func TestValidateModulePath(t *testing.T) {
	valid := []string{
		"github.com/myorg/payments",
		"example.com/svc",
		"github.com/MyOrg/payments",
	}
	for _, m := range valid {
		if err := ValidateModulePath(m); err != nil {
			t.Errorf("ValidateModulePath(%q) = %v, want nil", m, err)
		}
	}

	invalid := []string{
		"",
		"payments", // no dot in the first path element
		"github.com/my org/payments",
		"/github.com/myorg/payments",
		"github.com/myorg/payments/",
		"github.com/myorg//payments",
		"модуль.рф/payments",
	}
	for _, m := range invalid {
		if err := ValidateModulePath(m); err == nil {
			t.Errorf("ValidateModulePath(%q) = nil, want error", m)
		}
	}
}
