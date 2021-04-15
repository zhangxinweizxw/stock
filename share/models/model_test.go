package models

import (
    "database/sql"
    "fmt"
    "testing"

    "stock/share/gocraft/dbr"
)

type Roles struct {
	Model      int `db:"-"`
	IsAdmin    int
	IsDefault  int
	Money      float64
	Id         int64
	Email      dbr.NullString
	UpdateTime dbr.NullInt64
}

// StructToMap
func Test_StructToMap(t *testing.T) {
	Convey("StructToMap", t, func() {
		var r Roles
		r.Id = 100
		r.IsAdmin = 2
		r.Email = dbr.NullString{sql.NullString{"250135945@qq.com", true}}
		r.UpdateTime = dbr.NullInt64{sql.NullInt64{0, true}}
		m := StructToMap(r)

		So(m["Id"], ShouldEqual, int64(100))
		So(m["Email"], ShouldEqual, "250135945@qq.com")
		fmt.Printf("%v", m)
	})
}

// MapToStruct
func Test_MapToStruct(t *testing.T) {
	Convey("MapToStruct", t, func() {
		m := map[string]string{
			"IsAdmin":    "1",
			"IsDefault":  "2",
			"Money":      "10.11",
			"Id":         "100",
			"Email":      "aa",
			"UpdateTime": "0",
		}
		var r Roles
		err := MapToStruct(&r, m)
		So(err, ShouldBeNil)
		So(r.Id, ShouldEqual, int64(100))
		So(r.Email.String, ShouldEqual, "aa")
	})
}
