package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type tsignature struct {
	Id             int
	IsDelete       []byte
	Sms_content    string
	Sms_tsignature string
	Status         int
	Template_id    string
}

type Tsignatures struct {
	tsign []tsignature
}

type DbContext struct {
	db *sql.DB
}

func (self *DbContext) Open() *sql.DB {

	db, err := sql.Open("mysql", "root:121121@tcp(127.0.0.1:3306)/cussms?charset=utf8")
	if err != nil {
		Log(err)
		return nil
	}
	self.db = db
	return self.db
}

func (self *DbContext) Get(str string) string {
	if self.db == nil {
		Log("Dbcontext is nil")
	}
	t := tsignature{}
	row := self.db.QueryRow(`select * from tsignature where Sms_signature = '` + str + `'`)
	err := row.Scan(&t.Id, &t.IsDelete, &t.Sms_content, &t.Sms_tsignature, &t.Status, &t.Template_id)
	if err != nil {
		return ""
	}
	//Log(t)
	if t.Status != 2 {
		return ""
	}
	return t.Template_id
}

// func (self *DbContext) Remove() {
//
// }

func Log(v ...interface{}) {
	log.Println(v...)
}

func main() {

	dbcontext := DbContext{}
	dbcontext.Open()
	defer dbcontext.db.Close()

	str := dbcontext.Get("【Custouch】")
	if str == "" {
		Log("str is nil..")
	} else {
		Log(str)
	}

}

// m := IsDelete.(map[string]interface{})
// for k, v := range m {
//   switch vv := v.(type) {
//   case string:
//     Log(k, "is string", vv)
//   case int:
//     Log(k, "is int", vv)
//   case float64:
//     Log(k, "is float64", vv)
//   case []interface{}:
//     Log(k, "is an array:")
//     for i, u := range vv {
//       Log(i, u)
//     }
//   default:
//     Log(k, "is of a type I don't know how to handle")
//   }
// }
