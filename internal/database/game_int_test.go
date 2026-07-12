package database

import (
	"encoding/json"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestGameIntJSONUsesString(t *testing.T) {
	value := MustParseGameInt("1234567890123456789012345678901234567890")

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal game int: %v", err)
	}

	if string(data) != `"1234567890123456789012345678901234567890"` {
		t.Fatalf("unexpected JSON: %s", string(data))
	}

	var decoded GameInt
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal game int: %v", err)
	}

	if decoded.String() != value.String() {
		t.Fatalf("round-trip mismatch: got %s want %s", decoded.String(), value.String())
	}
}

func TestGameIntScansIntegerNumeric(t *testing.T) {
	var value GameInt

	if err := value.Scan(pgtype.Numeric{
		Int:   MustParseGameInt("123456789012345678901234567890").BigInt(),
		Exp:   0,
		Valid: true,
	}); err != nil {
		t.Fatalf("scan numeric: %v", err)
	}

	if value.String() != "123456789012345678901234567890" {
		t.Fatalf("scan mismatch: %s", value.String())
	}
}

func TestGameIntRejectsFractionalNumeric(t *testing.T) {
	var value GameInt

	err := value.Scan(pgtype.Numeric{
		Int:   MustParseGameInt("12345").BigInt(),
		Exp:   -2,
		Valid: true,
	})
	if err == nil {
		t.Fatal("expected fractional numeric scan to fail")
	}
}

func TestGameIntFormatsBase10(t *testing.T) {
	value := MustParseGameInt("123456789012345678901234567890")
	want := "123,456,789,012,345,678,901,234,567,890"

	if got := value.FormatBase10(); got != want {
		t.Fatalf("format mismatch: got %s want %s", got, want)
	}
}
