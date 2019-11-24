package obj

import (
	"reflect"
)

func Invoke(o interface{}, method string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	return reflect.ValueOf(o).MethodByName(method).Call(inputs)
}

func InvokeInt(o interface{}, method string, args ...interface{}) int {
	ret := Invoke(o, method, args...)
	return int(ret[0].Int())
}

func InvokeString(o interface{}, method string, args ...interface{}) string {
	ret := Invoke(o, method, args...)
	return ret[0].String()
}
