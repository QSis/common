package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/QSis/common/obj"
	. "github.com/QSis/exception"
	"github.com/jinzhu/gorm"
)

/********************************************************
 * 分页、排序、查询参数类型定义
 * 用于数据库查询、请求取参数
 ********************************************************/

type Pager struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
}

type Sorting struct {
	OrderBy string `json:"order_by"`
	Sort    string `json:"sort"`
}

type Query struct {
	Param      *map[string]interface{} `json:"param"`
	Column     *[]string               `json:"column"`
	Pager      *Pager                  `json:"pager"`
	Sorting    *Sorting                `json:"sorting"`
	FetchFirst bool                    `json:"fetch_first"`
}

var (
	IdDesc Sorting = Sorting{
		OrderBy: "id",
		Sort:    "desc",
	}

	IdAsc Sorting = Sorting{
		OrderBy: "id",
		Sort:    "asc",
	}
)

/********************************************************
 * 数据库查询方法：
 *   FindAll() 查询多条数据
 *   Find() 查询单条数据
 *
 * Param:
 *   data interface{}:  ORM 对象引用，例如 model.Script
 *   query interface{}: 查询 Query 对象，可以是 Query, Pager, Sorting, map[string]interface{}
 *     Pager 分页参数
 *     Sorting 排序参数
 *     Column 查询字段
 *     Query.Param, map[string]interface{} 查询参数
 *
 * Query.Param 示例:
 * map[string]interface{} {
 *   "id": 123,				// id = 123 等价操作符：eq, =, in，如："id:" : 123
 *	 "age:gt": 18,			// age > 18，等价操作符：gt, >，同类操作符：ge, >=, lt, <, le, <=
 *	 "gender:not": "M",		// gender != "M"，等价操作符：not, ne, !=, <>
 *	 "name:like": "wtf",	// name like "%wtf%"，等价操作符：like, %
 *	 "memo:prefix": "jd",	// memo like "jd%"，等价操作符：prefix, ^
 * }
 ********************************************************/

func FindAll(data interface{}, query ...interface{}) {
	db := &MyDB{DB.New()}
	db.FindAll(data, query...)
}

func (myDB *MyDB) FindAll(data interface{}, query ...interface{}) {
	db := myDB.DB
	db = parseQuery(db, query...)
	db = CheckResult(db.Find(data))
	var pager *Pager = nil
	for _, v := range query {
		switch v.(type) {
		case Query:
			q := v.(Query)
			pager = q.Pager
		case *Query:
			q := v.(*Query)
			pager = q.Pager
		case *Pager:
			pager = v.(*Pager)
		}
	}
	if pager != nil {
		var count int
		CheckResult(db.Offset(nil).Limit(nil).Count(&count))
		pager.TotalCount = count
	}
	// myDB.CopyDB(data)
}

func Count(data interface{}, query ...interface{}) int {
	db := &MyDB{DB.New()}
	return db.Count(data, query...)
}

func (myDB *MyDB) Count(data interface{}, query ...interface{}) int {
	db := myDB.DB
	var count int = 0
	db = parseQuery(db, query...)
	CheckResult(db.Model(data).Count(&count))
	return count
}

func Find(data interface{}, query ...interface{}) {
	db := &MyDB{DB.New()}
	db.Find(data, query...)
}

func (myDB *MyDB) Find(data interface{}, query ...interface{}) {
	db := myDB.DB
	db = parseQuery(db, query...)
	CheckResult(db.First(data))
	myDB.CopyDB(data)
}

func FindColumn(data interface{}, out interface{}, column string, query ...interface{}) {
	db := &MyDB{DB.New()}
	db.FindColumn(data, out, column, query...)
}

func (myDB *MyDB) FindColumn(data interface{}, out interface{}, column string, query ...interface{}) {
	db := myDB.DB
	db = parseQuery(db, query...)
	CheckResult(db.Model(data).Pluck(column, out))
}

func UpdateOne(data interface{}, to interface{}, query ...interface{}) int {
	db := &MyDB{DB.New()}
	return db.UpdateOne(data, to, query...)
}

func (myDB *MyDB) UpdateOne(data interface{}, to interface{}, query ...interface{}) int {
	query = append(query, Pager{
		Page:     1,
		PageSize: 1,
	})
	return myDB.UpdateAll(data, to, query...)
}

func UpdateAll(data interface{}, to interface{}, query ...interface{}) int {
	db := &MyDB{DB.New()}
	return db.UpdateAll(data, to, query...)
}

func (myDB *MyDB) UpdateAll(data interface{}, to interface{}, query ...interface{}) int {
	db := myDB.DB
	db = parseQuery(db, query...)
	db = CheckResult(db.Model(data).Update(to))
	return int(db.RowsAffected)
}

func DeleteAll(data interface{}, query ...interface{}) int {
	db := &MyDB{DB.New()}
	return db.DeleteAll(data, query...)
}

func (myDB *MyDB) DeleteAll(data interface{}, query ...interface{}) int {
	db := myDB.DB
	db = parseQuery(db, query...)
	db = CheckResult(db.Delete(data))
	return int(db.RowsAffected)
}

func (myDB *MyDB) CopyDB(data interface{}) {
	elem := reflect.ValueOf(data)
	if elem.Kind() == reflect.Ptr { // elem is a ptr type
		elem = elem.Elem()
	}
	if elem.Kind() == reflect.Struct {
		obj.Invoke(data, "SetDB", myDB)
	} else if elem.Kind() == reflect.Slice { // XXX
		/*
			elemLength := elem.Len()
			for i := 0; i < elemLength; i++ {
				fmt.Printf("I: %#v\n", elem.Index(i).Interface())
				obj.Invoke(elem.Index(i).Interface(), "SetDB", myDB)
			}
		*/
	}
}

/********************************************************************
 * Parse query condition
 ********************************************************************/

func (pager *Pager) Offset() int {
	return (pager.Page - 1) * pager.PageSize
}

func (pager *Pager) Load(db *gorm.DB) *gorm.DB {
	return db.Offset(pager.Offset()).Limit(pager.PageSize)
}

func (sorting *Sorting) Load(db *gorm.DB) *gorm.DB {
	return db.Order(sorting.OrderBy + " " + sorting.Sort)
}

func (query *Query) Load(db *gorm.DB) *gorm.DB {
	if query.Pager != nil {
		db = query.Pager.Load(db)
	}
	if query.Sorting != nil {
		db = query.Sorting.Load(db)
	}
	if query.Column != nil {
		db = db.Select(*query.Column)
	}
	if query.Param != nil {
		db = parseCondition(db, query.Param)
	}
	return db
}

func parseCondition(db *gorm.DB, param *map[string]interface{}) *gorm.DB {
	for key, v := range *param {
		keys := strings.Split(key, ":")
		op := "eq"
		k := keys[0]
		if len(keys) > 1 {
			op = strings.ToLower(keys[1])
		}
		switch op {
		case ">", "gt":
			db = db.Where(k+" > ?", v)
		case ">=", "ge":
			db = db.Where(k+" >= ?", v)
		case "<", "lt":
			db = db.Where(k+" < ?", v)
		case "<=", "le":
			db = db.Where(k+" <= ?", v)
		case "!=", "!", "ne", "not", "<>":
			switch v.(type) {
			case []string, []int:
				db = db.Where(k+" not in (?)", v)
			default:
				db = db.Where(k+" != ?", v)
			}
		case "like", "%":
			db = db.Where(k+" like ?", "%"+v.(string)+"%")
		case "prefix", "^":
			db = db.Where(k+" like ?", v.(string)+"%")
		case "=", "", "eq", "in":
			fallthrough
		default:
			switch v.(type) {
			case []string, []int:
				db = db.Where(k+" in (?)", v)
			default:
				db = db.Where(k+" = ?", v)
			}
		}
	}
	return db
}

func idCondition(db *gorm.DB, id interface{}) *gorm.DB {
	switch id.(type) {
	case int:
		db = db.Where([]int{id.(int)})
	case []int:
		db = db.Where(id.([]int))
	case string:
		db = db.Where([]string{id.(string)})
	case []string:
		db = db.Where(id.([]string))
	default:
		Throw(fmt.Sprintf("Invalid type %#v  in idCondition()", id))
	}
	return db
}

func parseQuery(db *gorm.DB, query ...interface{}) *gorm.DB {
	for _, v := range query {
		switch v.(type) {
		case *Query:
			db = v.(*Query).Load(db)
		case Query:
			q := v.(Query)
			db = q.Load(db)
		case *Pager:
			db = v.(*Pager).Load(db)
		case Pager:
			p := v.(Pager)
			db = p.Load(db)
		case *Sorting:
			db = v.(*Sorting).Load(db)
		case Sorting:
			s := v.(Sorting)
			db = s.Load(db)
		case map[string]interface{}:
			m := v.(map[string]interface{})
			db = parseCondition(db, &m)
		case *map[string]interface{}:
			db = parseCondition(db, v.(*map[string]interface{}))
		case int, []int, string, []string:
			db = idCondition(db, v)
		default:
			Throw(fmt.Sprintf("Invalid type %#v  in Query", v))
		}
	}
	return db
}
