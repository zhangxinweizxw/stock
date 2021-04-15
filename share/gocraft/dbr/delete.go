package dbr

import (
    "fmt"
	"strings"
)

// DeleteStmt builds `DELETE ...`
type DeleteStmt struct {
	raw

	Table string

	WhereCond []Condition
}

// Build builds `DELETE ...` in dialect
func (b *DeleteStmt) Build(d Dialect, buf Buffer) error {
	if b.raw.Query != "" {
		return b.raw.Build(d, buf)
	}

	if b.Table == "" {
		return ErrTableNotSpecified
	}

	schemaName := strings.Split(b.Table, ".")
	if len(schemaName) > 1 {
		buf.WriteString(fmt.Sprintf("/*!mycat:schema=%s*/", schemaName[0]))
	}

	buf.WriteString("DELETE FROM ")
	buf.WriteString(d.QuoteIdent(b.Table))

	if len(b.WhereCond) > 0 {
		buf.WriteString(" WHERE ")
		err := And(b.WhereCond...).Build(d, buf)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteFrom creates a DeleteStmt
func DeleteFrom(table string) *DeleteStmt {
	return &DeleteStmt{
		Table: table,
	}
}

// DeleteBySql creates a DeleteStmt from raw query
func DeleteBySql(query string, value ...interface{}) *DeleteStmt {
	return &DeleteStmt{
		raw: raw{
			Query: query,
			Value: value,
		},
	}
}

// Where adds a where condition
func (b *DeleteStmt) Where(query interface{}, value ...interface{}) *DeleteStmt {
	switch query := query.(type) {
	case string:
		b.WhereCond = append(b.WhereCond, Expr(query, value...))
	case Condition:
		b.WhereCond = append(b.WhereCond, query)
	}
	return b
}

// where adds
func (b *DeleteStmt) WhereMap(m map[string]interface{}) *DeleteStmt {
	for k, v := range m {
		b.WhereCond = append(b.WhereCond, ExprWhere(k, v))
	}
	return b
}
