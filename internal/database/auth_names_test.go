package database

import "testing"

func TestIsReservedAuthDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		reserved bool
	}{
		{name: "Soryn", reserved: true},
		{name: "soryn", reserved: true},
		{name: "S0ryn", reserved: true},
		{name: "S o r y n", reserved: true},
		{name: "Soryn!", reserved: true},
		{name: "Nami", reserved: true},
		{name: "Nami-chan", reserved: true},
		{name: "Nami chan", reserved: true},
		{name: "Nami_Chan", reserved: true},
		{name: "Namichan", reserved: true},
		{name: "N@mi-Chan", reserved: true},
		{name: "NamiFan", reserved: false},
		{name: "Sorin", reserved: false},
		{name: "Nana", reserved: false},
		{name: "MochaNami", reserved: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IsReservedAuthDisplayName(test.name)
			if got != test.reserved {
				t.Fatalf("expected reserved=%v, got %v", test.reserved, got)
			}
		})
	}
}

func TestReservedDisplayNameKey(t *testing.T) {
	tests := map[string]string{
		"Soryn":     "soryn",
		"S0ryn!":    "soryn",
		"Nami-chan": "namichan",
		"N@mi_Chan": "namichan",
		"Nami Fan":  "namifan",
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			got := reservedDisplayNameKey(input)
			if got != expected {
				t.Fatalf("expected %q, got %q", expected, got)
			}
		})
	}
}
