package types

import (
	"database/sql"
	"strings"
)

// ToNullString converts the given value to a nullable string
func ToNullString(value string) sql.NullString {
	value = strings.TrimSpace(value)
	return sql.NullString{
		Valid:  value != "",
		String: value,
	}
}
