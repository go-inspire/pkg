package log

import (
	"fmt"
	"reflect"
)

func sprint(a ...interface{}) string {
	if len(a) == 0 {
		return ""
	} else if s, ok := a[0].(string); ok && len(a) == 1 {
		return s
	} else if v := reflect.ValueOf(a[0]); len(a) == 1 && v.Kind() == reflect.String {
		return v.String()
	} else {
		return fmt.Sprint(a...)
	}
}
