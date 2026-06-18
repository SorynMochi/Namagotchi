package server

import (
	"os"
	"strings"
	"testing"
)

func TestSecuritySensitiveRouteWiring(t *testing.T) {
	sourceBytes, err := os.ReadFile("server.go")
	if err != nil {
		t.Fatalf("read server.go: %v", err)
	}

	source := string(sourceBytes)

	requiredSnippets := []string{
		`mux.HandleFunc("/api/auth/csrf", s.HandleCSRFToken)`,
		`mux.HandleFunc("/api/auth/logout", s.requireCSRF(s.HandleAuthLogout))`,
		`mux.HandleFunc("/api/dev/unlock", s.requireDev(s.requireCSRF(s.HandleDevUnlock)))`,
		`mux.HandleFunc("/api/dev/lock", s.requireDev(s.requireCSRF(s.HandleDevLock)))`,
		`mux.HandleFunc("/api/dev/force-tick", s.requireDev(s.requireDevUnlock(s.withDevAudit("force-tick", s.requireCSRF(s.HandleForceTick)))))`,
		`mux.HandleFunc("/api/dev/audit-logs", s.requireDev(s.requireDevUnlock(s.HandleDevAuditLogs)))`,
		`mux.HandleFunc("/api/player/sync", s.requireAuth(s.requireCSRF(s.HandlePlayerSync)))`,
		`mux.HandleFunc("/api/player/gathering", s.requireAuth(s.requireCSRF(s.HandleGatheringTask)))`,
		`mux.HandleFunc("/api/player/care", s.requireAuth(s.requireCSRF(s.HandleCareAction)))`,
		`mux.HandleFunc("/api/player/wardrobe/equip", s.requireAuth(s.requireCSRF(s.HandleEquipWardrobeItem)))`,
		`mux.HandleFunc("/api/player/wardrobe/unequip", s.requireAuth(s.requireCSRF(s.HandleUnequipWardrobeItem)))`,
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing security-sensitive route wiring: %s", snippet)
		}
	}
}
