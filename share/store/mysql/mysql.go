package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	MAX_RECONNECT_RETRIES = 5
)

var (
	db *sql.DB
)

// init db conn
func Init(driverName, dataSource string) error {
	if db == nil {
		var err error
		db, err = GetConn(driverName, dataSource)
		if err != nil {
			return err
		}
	}
	return nil
}

//get conn pool
func GetConn(driverName, dataSource string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSource)

	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(10)

	return db, err
}

// query mutil rows
func Query(cmd string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	for i := 0; i < MAX_RECONNECT_RETRIES; i++ {
		rows, err = db.Query(cmd, args...)
		if err == nil {
			break
		}
	}

	return rows, err
}

// query single row
func QueryRow(cmd string, args ...interface{}) *sql.Row {
	var row *sql.Row
	for i := 0; i < MAX_RECONNECT_RETRIES; i++ {
		row = db.QueryRow(cmd, args...)
		if row != nil {
			break
		}
	}
	return row
}

// exec
func Exec(cmd string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error
	for i := 0; i < MAX_RECONNECT_RETRIES; i++ {
		result, err = db.Exec(cmd, args...)
		if err == nil {
			break
		}
	}
	return result, err
}

// insert
func Insert(cmd string, args ...interface{}) (int64, error) {

	result, err := Exec(cmd, args...)

	if err != nil {
		return 0, err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

//update
func Update(cmd string, args ...interface{}) (int64, error) {
	result, err := Exec(cmd, args...)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	return count, nil
}

//delete
func Delete(cmd string, args ...interface{}) (int64, error) {
	result, err := Exec(cmd, args...)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	return count, err
}

//start transaction
func ExecTransaction() (*sql.Tx, error) {
	return db.Begin()
}

// close db pool
func Close() {
	if db != nil {
		db.Close()
	}
}
