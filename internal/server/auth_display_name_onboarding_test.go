package server

import (
	"testing"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

func TestDisplayNameChangeAllowedDuringOnboarding(t *testing.T) {
	placeholders := []string{
		"",
		"Player",
		"player",
		"Player_1",
		"Player_123",
		"NamiFan",
		"NamiFan2",
		"namifan99",
	}

	for _, placeholder := range placeholders {
		account := database.AuthAccount{DisplayName: placeholder}
		if !displayNameChangeAllowedDuringOnboarding(account, "FreshName") {
			t.Fatalf("expected placeholder %q to allow onboarding display name change", placeholder)
		}
	}
}

func TestDisplayNameChangeBlockedAfterOnboarding(t *testing.T) {
	account := database.AuthAccount{DisplayName: "Test_User1"}

	if displayNameChangeAllowedDuringOnboarding(account, "OtherName") {
		t.Fatal("expected custom display name to block free display name change")
	}
}

func TestDisplayNameChangeAllowsSameExistingName(t *testing.T) {
	account := database.AuthAccount{DisplayName: "Test_User1"}

	if !displayNameChangeAllowedDuringOnboarding(account, "test_user1") {
		t.Fatal("expected same existing display name to be allowed")
	}
}

func TestOnboardingPlaceholderRejectsLookalikeCustomNames(t *testing.T) {
	customNames := []string{
		"PlayerOne",
		"Player_X",
		"NamiFanGirl",
		"NamiFan_1",
	}

	for _, customName := range customNames {
		if isOnboardingPlaceholderDisplayName(customName) {
			t.Fatalf("expected %q not to count as onboarding placeholder", customName)
		}
	}
}
