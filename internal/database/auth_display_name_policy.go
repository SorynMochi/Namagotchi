package database

import "errors"

var ErrAuthDisplayNameInvalid = errors.New("auth display name is invalid")

const (
	AuthDisplayNameMinLength = 3
	AuthDisplayNameMaxLength = 15
)

func ValidateAuthDisplayName(displayName string) error {
	displayName = cleanAuthDisplayName(displayName)

	if len(displayName) < AuthDisplayNameMinLength || len(displayName) > AuthDisplayNameMaxLength {
		return ErrAuthDisplayNameInvalid
	}

	for _, r := range displayName {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '_' || r == '-' {
			continue
		}

		return ErrAuthDisplayNameInvalid
	}

	if IsReservedAuthDisplayName(displayName) {
		return ErrAuthDisplayNameReserved
	}

	return nil
}
