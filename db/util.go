package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/QSis/common/obj"
	. "github.com/QSis/exception"
	"github.com/jinzhu/gorm"
)

func Create(data interface{}) {
	(&MyDB{DB.New()}).Create(data)
}

func Save(data interface{}) {
	(&MyDB{DB.New()}).Save(data)
}

func Raw(sql string, binding ...interface{}) *gorm.DB {
	return (&MyDB{DB.New()}).Raw(sql, binding...)
}

func Exec(sql string, binding ...interface{}) *gorm.DB {
	return (&MyDB{DB.New()}).Exec(sql, binding...)
}

func Where(query interface{}, args ...interface{}) *gorm.DB {
	return (&MyDB{DB.New()}).Where(query, args...)
}

func BulkCreate(data interface{}) {
	(&MyDB{DB.New()}).BulkCreate(data)
}

// MyDB function

func (myDB *MyDB) Create(data interface{}) {
	CheckResult(myDB.DB.Create(data))
	myDB.CopyDB(data)
}

func (myDB *MyDB) Save(data interface{}) {
	CheckResult(myDB.DB.Save(data))
}

func (myDB *MyDB) Exec(sql string, binding ...interface{}) *gorm.DB {
	return CheckResult(myDB.DB.Exec(sql, binding...))
}

func (myDB *MyDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return myDB.DB.Where(query, args...)
}

func (myDB *MyDB) BulkCreate(data interface{}) {
	elem := reflect.ValueOf(data)
	if !elem.CanAddr() { // elem is a ptr type
		elem = elem.Elem()
	}
	elemLength := elem.Len()
	if elemLength <= 0 {
		Throw("Attempt to insert a empty slice to database")
	}
	scope := myDB.DB.NewScope(elem.Interface())

	column := []string{}
	fields := scope.Fields()
	for _, v := range fields {
		column = append(column, scope.Quote(v.DBName))
	}
	columnSql := strings.Join(column, ",")

	placeholder := "(" + obj.ImplodeRepeatString("?", ",", len(fields)) + ")"
	valueSql := obj.ImplodeRepeatString(placeholder, ",", elemLength)

	values := []interface{}{}
	for i := 0; i < elemLength; i++ {
		row := elem.Index(i)
		for _, field := range fields {
			v := row.FieldByName(field.Name).Interface()
			if field.IsPrimaryKey {
				switch v.(type) {
				case int:
					if v.(int) != 0 {
						values = append(values, v.(int))
					} else {
						values = append(values, nil)
					}
				case int64:
					if v.(int64) != 0 {
						values = append(values, v.(int64))
					} else {
						values = append(values, nil)
					}
				default:
					values = append(values, nil)
				}
			} else {
				values = append(values, v)
			}
		}
	}

	table := scope.QuotedTableName()
	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES %s", table, columnSql, valueSql)
	result, err := scope.SQLDB().Exec(sql, values...)
	if err != nil {
		Throw(err)
	}
	primaryKeyField := scope.PrimaryField()
	if primaryKeyField == nil {
		return
	}
	lastId, err2 := result.LastInsertId()
	if err2 != nil {
		Throw(err2)
	}

	for i := 0; i < elemLength; i++ {
		row := elem.Index(i)
		row.FieldByName(primaryKeyField.Name).SetInt(lastId)
		lastId++
	}
	// myDB.CopyDB(data)
}

func CheckResult(db *gorm.DB) *gorm.DB {
	if db.Error != nil {
		Throw(db.Error)
	}
	return db
}

func ExceptionIsRecordNotFound(e interface{}) bool {
	switch e.(type) {
	case error:
		return e.(error) == ErrRecordNotFound
	}
	return false
}
