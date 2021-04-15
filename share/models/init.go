package models

import (
    "sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"stock/share/gocraft/dbr"
)

var (
	MyCat *dbr.Session
	l     sync.RWMutex
)

func Init(driverName, dataSource string) error {
	l.Lock()
	defer l.Unlock()

	if MyCat == nil {
		conn, err := dbr.Open(
			driverName,
			dataSource, nil)

		if err != nil {
			return err
		}

		conn.SetMaxOpenConns(50)
		conn.SetMaxIdleConns(10)
		conn.SetConnMaxLifetime(time.Second * 100)

		MyCat = conn.NewSession(nil)
		if MyCat == nil {
			return err
		}

	}
	return nil
}
