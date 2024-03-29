package dbr

import (
    "fmt"
)

// XxxBuilders all support raw query
type raw struct {
	Query string
	Value []interface{}
}

// Expr should be used when sql syntax is not supported
func Expr(query string, value ...interface{}) Builder {
	return &raw{Query: query, Value: value}
}

// Expr should be used when sql syntax is not supported
func ExprWhere(query string, value ...interface{}) Builder {
	return &raw{Query: fmt.Sprintf("%s=?", query), Value: value}
}

func (raw *raw) Build(_ Dialect, buf Buffer) error {
	buf.WriteString(raw.Query)
	buf.WriteValue(raw.Value...)
	return nil
}
