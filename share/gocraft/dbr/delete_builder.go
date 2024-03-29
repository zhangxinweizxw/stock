package dbr

import (
	"database/sql"
    "fmt"
)

type DeleteBuilder struct {
	runner
	EventReceiver
	Dialect Dialect

	*DeleteStmt

	LimitCount int64
}

func (sess *Session) DeleteFrom(table string) *DeleteBuilder {
	return &DeleteBuilder{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		DeleteStmt:    DeleteFrom(table),
		LimitCount:    -1,
	}
}

func (tx *Tx) DeleteFrom(table string) *DeleteBuilder {
	return &DeleteBuilder{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		DeleteStmt:    DeleteFrom(table),
		LimitCount:    -1,
	}
}

func (sess *Session) DeleteBySql(query string, value ...interface{}) *DeleteBuilder {
	return &DeleteBuilder{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		DeleteStmt:    DeleteBySql(query, value...),
		LimitCount:    -1,
	}
}

func (tx *Tx) DeleteBySql(query string, value ...interface{}) *DeleteBuilder {
	return &DeleteBuilder{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		DeleteStmt:    DeleteBySql(query, value...),
		LimitCount:    -1,
	}
}

func (b *DeleteBuilder) ToSql() (string, []interface{}) {
	buf := NewBuffer()
	err := b.Build(b.Dialect, buf)
	if err != nil {
		panic(err)
	}
	return buf.String(), buf.Value()
}

func (b *DeleteBuilder) Exec() (sql.Result, error) {
	return exec(b.runner, b.EventReceiver, b, b.Dialect)
}

func (b *DeleteBuilder) Where(query interface{}, value ...interface{}) *DeleteBuilder {
	b.DeleteStmt.Where(query, value...)
	return b
}

func (b *DeleteBuilder) WhereMap(m map[string]interface{}) *DeleteBuilder {
	b.DeleteStmt.WhereMap(m)
	return b
}

func (b *DeleteBuilder) Limit(n uint64) *DeleteBuilder {
	b.LimitCount = int64(n)
	return b
}

func (b *DeleteBuilder) Build(d Dialect, buf Buffer) error {
	err := b.DeleteStmt.Build(b.Dialect, buf)
	if err != nil {
		return err
	}
	if b.LimitCount >= 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(fmt.Sprint(b.LimitCount))
	}
	return nil
}
