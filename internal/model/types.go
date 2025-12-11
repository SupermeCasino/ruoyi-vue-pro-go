package model

import (
	"database/sql/driver"
	"errors"
)

// BitBool is a boolean that maps to BIT(1) in database
type BitBool bool

// Scan implements the Scanner interface.
func (b *BitBool) Scan(value interface{}) error {
	if value == nil {
		*b = false
		return nil
	}

	switch v := value.(type) {
	case []uint8:
		if len(v) > 0 {
			*b = BitBool(v[0] == 1)
		} else {
			*b = false
		}
	case int64:
		*b = BitBool(v == 1)
	case bool:
		*b = BitBool(v)
	default:
		return errors.New("incompatible type for BitBool")
	}
	return nil
}

// Value implements the driver Valuer interface.
func (b BitBool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

func NewBitBool(b bool) BitBool {
	return BitBool(b)
}
