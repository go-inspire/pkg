/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"fmt"
)

// sprint 将任意类型参数转换为字符串
// 功能说明:
// - 如果没有参数，返回空字符串
// - 如果只有一个参数:
//   - 如果是字符串类型，直接返回
//   - 如果实现了 fmt.Stringer 接口，调用其 String() 方法
//   - 其他情况使用 fmt.Sprint 转换
// - 多个参数时直接使用 fmt.Sprint 转换
func sprint(a ...interface{}) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		// 处理nil值
		if a[0] == nil {
			return "<nil>"
		}
		
		// 使用类型switch替代多次类型断言，更高效
		switch v := a[0].(type) {
		case string:
			return v
		case fmt.Stringer:
			return v.String()
		default:
			return fmt.Sprint(v)
		}
	default:
		return fmt.Sprint(a...)
	}
}

// sprintf 格式化字符串模板
// 功能说明:
// - 没有参数时返回原模板
// - 模板非空时使用 fmt.Sprintf 格式化
// - 只有一个字符串参数且模板为空时直接返回该字符串
// - 其他情况调用 sprint 转换参数
func sprintf(template string, args ...interface{}) string {
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
	return sprint(args...)
}
