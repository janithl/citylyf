package entities_test

import (
	"testing"

	"github.com/janithl/citylyf/internal/entities"
)

func TestGetPlaceName(t *testing.T) {
	ns := entities.NewNameService()
	place := ns.GetPlaceName()
	if place == "" {
		t.Errorf(`GetPlaceName() = %q, got "", wanted a non-empty string, error`, place)
	}
}

func TestGetCompanyName(t *testing.T) {
	ns := entities.NewNameService()
	company := ns.GetCompanyName()
	if company == "" {
		t.Errorf(`GetCompanyName() = %q, got "", wanted a non-empty string, error`, company)
	}
}
