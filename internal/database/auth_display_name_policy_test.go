package database

import (
	"errors"
	"testing"
)

func TestValidateAuthDisplayNameAllowsSafeNames(t *testing.T) {
	validNames := []string{
		"Pixie_77",
		"Nami-Fan1",
		"abc",
		"ABC_123-xyz",
		"Player-001",
	}

	for _, name := range validNames {
		if err := ValidateAuthDisplayName(name); err != nil {
			t.Fatalf("expected %q to be valid, got %v", name, err)
		}
	}
}

func TestValidateAuthDisplayNameRejectsInvalidNames(t *testing.T) {
	invalidNames := []string{
		"",
		"ab",
		"this_name_is_too_long",
		"has space",
		"has.dot",
		"has@symbol",
		"emoji_ðŸ’–",
	}

	for _, name := range invalidNames {
		err := ValidateAuthDisplayName(name)
		if !errors.Is(err, ErrAuthDisplayNameInvalid) {
			t.Fatalf("expected %q to be invalid, got %v", name, err)
		}
	}
}

func TestValidateAuthDisplayNameRejectsReservedNames(t *testing.T) {
	reservedNames := []string{
		"Soryn",
		"soryn",
		"s0ryn",
		"Nami",
		"Nami-chan",
		"NAMI_CHAN",
		"Namichan",
	}

	for _, name := range reservedNames {
		err := ValidateAuthDisplayName(name)
		if !errors.Is(err, ErrAuthDisplayNameReserved) {
			t.Fatalf("expected %q to be reserved, got %v", name, err)
		}
	}
}
