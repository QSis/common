package db

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/QSis/common/obj"
)

type Field struct {
	Name       string
	Type       string
	Len        int
	PrimaryKey bool
	Key        []string
	Unique     bool
	Comment    string
}

func (f Field) DbColumnType() string {
	if len(f.Key) > 0 || f.Type == "" {
		return ""
	}
	if f.Type == "int" {
		return "int"
	} else if f.Type == "string" && f.Len == 0 {
		return "text"
	}
	return "varchar"
}

func (f Field) Id() string {
	if len(f.Key) > 0 {
		return f.Name
	}
	return "__" + f.Name
}

func (f Field) EqualTo(f2 Field) bool {
	return (f.Name == f2.Name && f.DbColumnType() == f2.DbColumnType() && f.Len == f2.Len)
}

func (f Field) SubSql() string {
	if len(f.Key) > 0 {
		t := "KEY"
		if f.Unique {
			t = "UNIQUE KEY"
		}
		keys := f.Key
		for k, v := range keys {
			keys[k] = "`" + v + "`"
		}
		return fmt.Sprintf("%s(%s)", t, strings.Join(keys, ", "))
	}

	sql := ""
	switch f.Type {
	case "bigint":
		fallthrough
	case "int":
		sql = fmt.Sprintf("`%s` %s DEFAULT 0 COMMENT '%s'", f.Name, f.Type, f.Comment)
	case "string":
		if f.Len > 0 {
			sql = fmt.Sprintf("`%s` varchar(%d) DEFAULT '' COMMENT '%s'", f.Name, f.Len, f.Comment)
		} else {
			sql = fmt.Sprintf("`%s` text DEFAULT '' COMMENT '%s'", f.Name, f.Comment)
		}
	}
	return sql
}

// doDrop: DEFAULT false
func GenerateAlterSql(data interface{}, doDrop ...bool) string {
	actual := ScanTable(data)
	if len(actual) == 0 {
		return GenerateCreateSql(data)
	}

	t := reflect.TypeOf(data)
	expect := []Field{}
	ScanStruct(&expect, t)

	mExpect := map[string]Field{}
	mActual := map[string]Field{}
	for _, v := range expect {
		mExpect[v.Id()] = v
	}
	for _, v := range actual {
		mActual[v.Id()] = v
	}

	// do diff
	subSql := []string{}
	for _, v := range expect {
		id := v.Id()
		v2, ok := mActual[id]
		if !ok {
			subSql = append(subSql, "ADD "+v.SubSql())
		} else if !v.EqualTo(v2) {
			subSql = append(subSql, "MODIFY "+v.SubSql())
		}
	}

	if len(doDrop) > 0 && doDrop[0] {
		for _, v := range actual {
			id := v.Id()
			if _, ok := mActual[id]; !ok {
				subSql = append(subSql, "DROP "+v.Name)
			}
		}
	}

	if len(subSql) == 0 {
		return ""
	}
	return fmt.Sprintf("ALTER TABLE %s %s", DB.NewScope(data).QuotedTableName(), strings.Join(subSql, ", "))
}

func ScanTable(data interface{}) []Field {
	fields := []Field{}
	columns := []struct {
		Field string `gorm:"column:Field"`
		Type  string `gorm:"column:Type"`
	}{}
	DB.New().Raw("desc " + DB.NewScope(data).QuotedTableName()).Scan(&columns)
	if len(columns) == 0 {
		return fields
	}
	for _, v := range columns {
		t := v.Type
		l := 0
		if strings.Index(t, "(") > 0 {
			t = t[:strings.Index(t, "(")]
			if t == "varchar" {
				l = obj.ToInt(strings.Trim(v.Type, "varchar()"))
			}
		}
		fields = append(fields, Field{
			Name: v.Field,
			Type: t,
			Len:  l,
		})
	}

	keys := []struct {
		Name    string `gorm:"column:Key_name"`
		Column  string `gorm:"column:Column_name"`
		NonUniq int    `gorm:"column:Non_unique"`
	}{}
	DB.New().Raw("show keys in " + DB.NewScope(data).QuotedTableName() + " where Key_name != 'PRIMARY'").Scan(&keys)
	keyMap := map[string]*Field{}
	for _, v := range keys {
		if _, ok := keyMap[v.Name]; !ok {
			keyMap[v.Name] = &Field{
				Key:    []string{},
				Unique: v.NonUniq == 0,
			}
		}
		keyMap[v.Name].Key = append(keyMap[v.Name].Key, v.Column)
	}

	for _, v := range keyMap {
		sort.Strings(v.Key)
		v.Name = strings.Join(v.Key, "__")
		fields = append(fields, *v)
	}
	return fields
}

func GenerateCreateSql(data interface{}) string {
	sql := []string{"CREATE TABLE " + DB.NewScope(data).QuotedTableName() + "("}
	fieldSql := []string{}
	fields := []Field{}
	t := reflect.TypeOf(data)
	ScanStruct(&fields, t)
	for _, v := range fields {
		switch v.Type {
		case "bigint":
			fallthrough
		case "int":
			if v.PrimaryKey {
				fieldSql = append(fieldSql, fmt.Sprintf("    `%s` %s AUTO_INCREMENT COMMENT '%s'", v.Name, v.Type, v.Comment))
			} else {
				fieldSql = append(fieldSql, fmt.Sprintf("    `%s` %s DEFAULT 0 COMMENT '%s'", v.Name, v.Type, v.Comment))
			}
		case "string":
			if v.Len > 0 {
				fieldSql = append(fieldSql, fmt.Sprintf("    `%s` varchar(%d) DEFAULT '' COMMENT '%s'", v.Name, v.Len, v.Comment))
			} else {
				fieldSql = append(fieldSql, fmt.Sprintf("    `%s` text DEFAULT '' COMMENT '%s'", v.Name, v.Comment))
			}
		}
	}
	for _, v := range fields {
		if v.PrimaryKey {
			fieldSql = append(fieldSql, fmt.Sprintf("    PRIMARY KEY `idx_%s` (`%s`)", v.Name, v.Name))
		}
	}
	for _, v := range fields {
		if len(v.Key) > 0 {
			keyQuoted := []string{}
			for _, k := range v.Key {
				keyQuoted = append(keyQuoted, "`"+k+"`")
			}
			keyType := "KEY"
			indexPrefix := "idx"
			if v.Unique {
				keyType = "UNIQUE KEY"
				indexPrefix = "uniq"
			}
			fieldSql = append(fieldSql, fmt.Sprintf("    %s `%s_%s` (%s)", keyType, indexPrefix, strings.Join(v.Key, "_"), strings.Join(keyQuoted, ",")))
		}
	}
	sql = append(sql, strings.Join(fieldSql, ",\n"))
	sql = append(sql, fmt.Sprintf(") ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='%s';", DB.NewScope(data).TableName()))
	return strings.Join(sql, "\n")
}

func ScanStruct(out *[]Field, t reflect.Type) {
	keyField := []Field{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Tag.Get("json")
		gormTag := f.Tag.Get("gorm")
		if gormTag == "-" {
			continue
		}
		field := Field{
			Name:       name,
			PrimaryKey: (gormTag == "primary_key"),
			Type:       "",
			Comment:    f.Tag.Get("comment"),
		}
		switch f.Type.Kind() {
		case reflect.Int:
			field.Type = "int"
		case reflect.Int64:
			field.Type = "bigint"
		case reflect.String:
			field.Type = "string"
			len := f.Tag.Get("len")
			if len == "text" {
				field.Len = -1
			} else {
				field.Len = obj.ToInt(len)
			}
			if field.Len == 0 {
				field.Len = 127
			}
		case reflect.Struct:
			ScanStruct(out, f.Type)
		}
		if field.Type != "" {
			*out = append(*out, field)
		}

		if val, ok := f.Tag.Lookup("unique"); ok {
			keyField = append(keyField, GetKeyField(val, name, true))
		}
		if val, ok := f.Tag.Lookup("key"); ok {
			keyField = append(keyField, GetKeyField(val, name, false))
		}
	}
	*out = append(*out, keyField...)
}

func GetKeyField(tag string, columnName string, unique bool) Field {
	keyMap := map[string]bool{
		columnName: true,
	}
	if tag != "" {
		seq := strings.Split(tag, ",")
		for _, v := range seq {
			keyMap[v] = true
		}
	}

	key := []string{}
	for k, _ := range keyMap {
		key = append(key, k)
	}
	sort.Strings(key)
	return Field{
		Name:   strings.Join(key, "__"),
		Key:    key,
		Unique: unique,
	}
}
