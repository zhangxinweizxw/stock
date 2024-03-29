package dbr

import (
    "fmt"
	"strings"
)

// SelectStmt builds `SELECT ...`
type SelectStmt struct {
	raw

	IsDistinct  bool
	annontation string

	Column    []interface{}
	Table     interface{}
	JoinTable []Builder

	WhereCond  []Condition
	Group      []Builder
	HavingCond []Condition
	Order      []Builder

	LimitCount  int64
	OffsetCount int64
}

// Build builds `SELECT ...` in dialect
func (b *SelectStmt) Build(d Dialect, buf Buffer) error {
	if b.raw.Query != "" {
		return b.raw.Build(d, buf)
	}

	if len(b.Column) == 0 {
		return ErrColumnNotSpecified
	}

	table := b.Table.(string)
	schemaName := strings.Split(table, ".")
	if len(schemaName) > 1 {
		buf.WriteString(fmt.Sprintf("/*!mycat:schema=%s*/", schemaName[0]))
	}

	buf.WriteString("SELECT ")

	if b.IsDistinct {
		buf.WriteString("DISTINCT ")
	}

	for i, col := range b.Column {
		if i > 0 {
			buf.WriteString(", ")
		}
		switch col := col.(type) {
		case string:
			// FIXME: no quote ident
			buf.WriteString(col)
		default:
			buf.WriteString(d.Placeholder())
			buf.WriteValue(col)
		}
	}

	if b.Table != nil {
		buf.WriteString(" FROM ")
		switch table := b.Table.(type) {
		case string:
			// FIXME: no quote ident
			buf.WriteString(table)
		default:
			buf.WriteString(d.Placeholder())
			buf.WriteValue(table)
		}
		if len(b.JoinTable) > 0 {
			for _, join := range b.JoinTable {
				err := join.Build(d, buf)
				if err != nil {
					return err
				}
			}
		}
	}

	if len(b.WhereCond) > 0 {
		buf.WriteString(" WHERE ")
		err := And(b.WhereCond...).Build(d, buf)
		if err != nil {
			return err
		}
	}

	if len(b.Group) > 0 {
		buf.WriteString(" GROUP BY ")
		for i, group := range b.Group {
			if i > 0 {
				buf.WriteString(", ")
			}
			err := group.Build(d, buf)
			if err != nil {
				return err
			}
		}
	}

	if len(b.HavingCond) > 0 {
		buf.WriteString(" HAVING ")
		err := And(b.HavingCond...).Build(d, buf)
		if err != nil {
			return err
		}
	}

	if len(b.Order) > 0 {
		buf.WriteString(" ORDER BY ")
		for i, order := range b.Order {
			if i > 0 {
				buf.WriteString(", ")
			}
			err := order.Build(d, buf)
			if err != nil {
				return err
			}
		}
	}

	if b.LimitCount >= 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(fmt.Sprint(b.LimitCount))
	}

	if b.OffsetCount >= 0 {
		buf.WriteString(" OFFSET ")
		buf.WriteString(fmt.Sprint(b.OffsetCount))
	}
	return nil
}

// Select creates a SelectStmt
func Select(column ...interface{}) *SelectStmt {
	return &SelectStmt{
		Column:      column,
		LimitCount:  -1,
		OffsetCount: -1,
	}
}

// From specifies table
func (b *SelectStmt) From(table interface{}) *SelectStmt {
	b.Table = table
	return b
}

// SelectBySql creates a SelectStmt from raw query
func SelectBySql(query string, value ...interface{}) *SelectStmt {
	return &SelectStmt{
		raw: raw{
			Query: query,
			Value: value,
		},
		LimitCount:  -1,
		OffsetCount: -1,
	}
}

// Distinct adds `DISTINCT`
func (b *SelectStmt) Distinct() *SelectStmt {
	b.IsDistinct = true
	return b
}

// Where adds a where condition
func (b *SelectStmt) Where(query interface{}, value ...interface{}) *SelectStmt {
	switch query := query.(type) {
	case string:
		b.WhereCond = append(b.WhereCond, Expr(query, value...))
	case Condition:
		b.WhereCond = append(b.WhereCond, query)
	}
	return b
}

func (b *SelectStmt) WhereMap(m map[string]interface{}) *SelectStmt {
	for k, v := range m {
		b.WhereCond = append(b.WhereCond, ExprWhere(k, v))
	}
	return b
}

// Having adds a having condition
func (b *SelectStmt) Having(query interface{}, value ...interface{}) *SelectStmt {
	switch query := query.(type) {
	case string:
		b.HavingCond = append(b.HavingCond, ExprWhere(query, value...))
	case Condition:
		b.HavingCond = append(b.HavingCond, query)
	}
	return b
}

// GroupBy specifies columns for grouping
func (b *SelectStmt) GroupBy(col ...string) *SelectStmt {
	for _, group := range col {
		b.Group = append(b.Group, Expr(group))
	}
	return b
}

// OrderBy specifies columns for ordering
func (b *SelectStmt) OrderAsc(col string) *SelectStmt {
	b.Order = append(b.Order, Order(col, ASC))
	return b
}

func (b *SelectStmt) OrderDesc(col string) *SelectStmt {
	b.Order = append(b.Order, Order(col, DESC))
	return b
}

// Limit adds limit
func (b *SelectStmt) Limit(n uint64) *SelectStmt {
	b.LimitCount = int64(n)
	return b
}

// Offset adds offset
func (b *SelectStmt) Offset(n uint64) *SelectStmt {
	b.OffsetCount = int64(n)
	return b
}

// Join joins table on condition
func (b *SelectStmt) Join(table, on interface{}) *SelectStmt {
	b.JoinTable = append(b.JoinTable, Join(Inner, table, on))
	return b
}

func (b *SelectStmt) LeftJoin(table, on interface{}) *SelectStmt {
	b.JoinTable = append(b.JoinTable, Join(Left, table, on))
	return b
}

func (b *SelectStmt) RightJoin(table, on interface{}) *SelectStmt {
	b.JoinTable = append(b.JoinTable, Join(Right, table, on))
	return b
}

func (b *SelectStmt) FullJoin(table, on interface{}) *SelectStmt {
	b.JoinTable = append(b.JoinTable, Join(Full, table, on))
	return b
}

// As creates alias for select statement
func (b *SelectStmt) As(alias string) Builder {
	return as(b, alias)
}
