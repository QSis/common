package obj

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/remrain/weakjson"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func PrettyPrint(data interface{}) {
	fmt.Println(PrettyEncode(data))
}

func PrettyEncode(data interface{}) string {
	encoded, _ := weakjson.MarshalIndent(data, "", "  ")
	return string(encoded)
}

func Encode(data interface{}) string {
	encoded, _ := weakjson.Marshal(data)
	return string(encoded)
}

func Decode(data interface{}, encoded string) error {
	return weakjson.Unmarshal([]byte(encoded), data)
}

func Copy(to interface{}, from interface{}) error {
	encoded, err := weakjson.Marshal(from)
	if err != nil {
		return err
	}
	err = weakjson.Unmarshal(encoded, to)
	return err
}
func ReflectCopy(dst interface{}, src interface{}) (err error) {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		err = errors.New("dst isn't a pointer to struct")
		return
	}
	dstElem := dstValue.Elem()
	if dstElem.Kind() != reflect.Struct {
		err = errors.New("pointer doesn't point to struct")
		return
	}

	srcValue := reflect.ValueOf(src)
	srcType := reflect.TypeOf(src)
	if srcType.Kind() != reflect.Struct {
		err = errors.New("src isn't struct")
		return
	}

	for i := 0; i < srcType.NumField(); i++ {
		sf := srcType.Field(i)
		sv := srcValue.FieldByName(sf.Name)
		// make sure the value which in dst is valid and can set
		if dv := dstElem.FieldByName(sf.Name); dv.IsValid() && dv.CanSet() {
			dv.Set(sv)
		}
	}
	return
}

func Contains(s interface{}, elem interface{}) bool {
	v := reflect.ValueOf(s)
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {

			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if v.Index(i).Interface() == elem {
				return true
			}
		}
	case reflect.String:
		return strings.Contains(s.(string), elem.(string))
	default:
		return false
	}
	return false
}

func ToString(v int) string {
	return strconv.Itoa(v)
}

func ToInt(v string) int {
	value, _ := strconv.Atoi(v)
	return value
}

func Retry(desc string, retryTimes int, method func() error) (err error) {
	for i := 1; i < retryTimes+1; i++ {
		err = method()
		if err == nil {
			return
		} else if err == sql.ErrNoRows {
			seelog.Warnf("%s no rows found: %s", desc, err)
			return
		}
		seelog.Warnf("%s failed: %s, retry %d times", desc, err, i)
		time.Sleep(time.Duration(i*i) * time.Second)
	}
	seelog.Warnf("%s failed after retry %d times, last error: %s", desc, retryTimes, err)
	return err
}
