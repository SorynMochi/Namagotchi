package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// GameInt is an exact integer for idle-game quantities that can grow beyond int64.
//
// It marshals to JSON as a string so browser code can safely parse it with BigInt
// instead of losing precision through JavaScript's Number type.
type GameInt struct {
	value *big.Int
}

func NewGameInt(value *big.Int) GameInt {
	if value == nil {
		return ZeroGameInt()
	}

	return GameInt{value: new(big.Int).Set(value)}
}

func ZeroGameInt() GameInt {
	return GameInt{value: new(big.Int)}
}

func GameIntFromInt64(value int64) GameInt {
	return GameInt{value: big.NewInt(value)}
}

func MustParseGameInt(value string) GameInt {
	parsed, err := ParseGameInt(value)
	if err != nil {
		panic(err)
	}

	return parsed
}

func ParseGameInt(value string) (GameInt, error) {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return ZeroGameInt(), nil
	}

	parsed, ok := new(big.Int).SetString(cleaned, 10)
	if !ok {
		return ZeroGameInt(), fmt.Errorf("invalid game integer: %q", value)
	}

	return NewGameInt(parsed), nil
}

func (n GameInt) BigInt() *big.Int {
	if n.value == nil {
		return new(big.Int)
	}

	return new(big.Int).Set(n.value)
}

func (n GameInt) String() string {
	if n.value == nil {
		return "0"
	}

	return n.value.String()
}

func (n GameInt) IsZero() bool {
	return n.value == nil || n.value.Sign() == 0
}

func (n GameInt) Sign() int {
	if n.value == nil {
		return 0
	}

	return n.value.Sign()
}

func (n GameInt) Cmp(other GameInt) int {
	return n.BigInt().Cmp(other.BigInt())
}

func (n GameInt) Add(other GameInt) GameInt {
	return NewGameInt(new(big.Int).Add(n.BigInt(), other.BigInt()))
}

func (n GameInt) Sub(other GameInt) GameInt {
	return NewGameInt(new(big.Int).Sub(n.BigInt(), other.BigInt()))
}

func (n GameInt) Mul(other GameInt) GameInt {
	return NewGameInt(new(big.Int).Mul(n.BigInt(), other.BigInt()))
}

func (n GameInt) Div(other GameInt) GameInt {
	if other.Sign() == 0 {
		return ZeroGameInt()
	}

	return NewGameInt(new(big.Int).Div(n.BigInt(), other.BigInt()))
}

func (n GameInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

func (n *GameInt) UnmarshalJSON(data []byte) error {
	if n == nil {
		return fmt.Errorf("cannot unmarshal game integer into nil receiver")
	}

	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		*n = ZeroGameInt()
		return nil
	}

	var text string
	if strings.HasPrefix(trimmed, "\"") {
		if err := json.Unmarshal(data, &text); err != nil {
			return err
		}
	} else {
		text = trimmed
	}

	parsed, err := ParseGameInt(text)
	if err != nil {
		return err
	}

	*n = parsed
	return nil
}

func (n *GameInt) Scan(src any) error {
	if n == nil {
		return fmt.Errorf("cannot scan game integer into nil receiver")
	}

	switch value := src.(type) {
	case nil:
		*n = ZeroGameInt()
		return nil
	case int64:
		*n = GameIntFromInt64(value)
		return nil
	case int32:
		*n = GameIntFromInt64(int64(value))
		return nil
	case int:
		*n = GameIntFromInt64(int64(value))
		return nil
	case string:
		return n.scanString(value)
	case []byte:
		return n.scanString(string(value))
	case pgtype.Numeric:
		return n.scanNumeric(value)
	default:
		return fmt.Errorf("unsupported game integer scan type %T", src)
	}
}

func (n *GameInt) scanString(value string) error {
	parsed, err := ParseGameInt(value)
	if err != nil {
		return err
	}

	*n = parsed
	return nil
}

func (n *GameInt) scanNumeric(value pgtype.Numeric) error {
	if !value.Valid {
		*n = ZeroGameInt()
		return nil
	}

	if value.NaN || value.InfinityModifier != pgtype.Finite {
		return fmt.Errorf("cannot scan non-finite numeric as game integer")
	}

	integer := new(big.Int).Set(value.Int)

	if value.Exp > 0 {
		integer.Mul(integer, pow10(int(value.Exp)))
	} else if value.Exp < 0 {
		divisor := pow10(int(-value.Exp))
		quotient, remainder := new(big.Int), new(big.Int)
		quotient.QuoRem(integer, divisor, remainder)
		if remainder.Sign() != 0 {
			return fmt.Errorf("cannot scan fractional numeric as game integer: %s", value.Int.String())
		}
		integer = quotient
	}

	*n = NewGameInt(integer)
	return nil
}

func (n GameInt) Value() (driver.Value, error) {
	return n.String(), nil
}

func (n GameInt) Int64() (int64, error) {
	if !n.BigInt().IsInt64() {
		return 0, fmt.Errorf("game integer %s overflows int64", n.String())
	}

	return n.BigInt().Int64(), nil
}

func (n GameInt) MustInt64() int64 {
	value, err := n.Int64()
	if err != nil {
		panic(err)
	}

	return value
}

func (n GameInt) FormatBase10() string {
	return addDigitGrouping(n.String())
}

func addDigitGrouping(value string) string {
	if value == "" {
		return "0"
	}

	sign := ""
	if strings.HasPrefix(value, "-") {
		sign = "-"
		value = strings.TrimPrefix(value, "-")
	}

	if len(value) <= 3 {
		return sign + value
	}

	firstGroupLength := len(value) % 3
	if firstGroupLength == 0 {
		firstGroupLength = 3
	}

	var builder strings.Builder
	builder.WriteString(sign)
	builder.WriteString(value[:firstGroupLength])

	for index := firstGroupLength; index < len(value); index += 3 {
		builder.WriteString(",")
		builder.WriteString(value[index : index+3])
	}

	return builder.String()
}

func pow10(exp int) *big.Int {
	if exp <= 0 {
		return big.NewInt(1)
	}

	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exp)), nil)
}

func (n GameInt) GoString() string {
	return "GameInt(" + strconv.Quote(n.String()) + ")"
}
