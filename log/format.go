/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"fmt"
	"reflect"
)

func Sprint(a ...interface{}) string {
	if len(a) == 0 {
		return ""
	} else if len(a) == 1 {
		if s, ok := a[0].(string); ok {
			return s
		} else if v := reflect.ValueOf(a[0]); v.Kind() == reflect.String {
			return v.String()
		} else {
			return fmt.Sprint(a...)
		}
	} else {
		return fmt.Sprint(a...)
	}
}

func Sprintf(template string, args ...interface{}) string {
	if len(args) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, args...)
	}

	if len(args) == 1 {
		if str, ok := args[0].(string); ok {
			return str
		}
	}
	return Sprint(args...)
}
