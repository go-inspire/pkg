/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"context"
	"fmt"
	"strconv"
)

func ExampleNewSharedChannel() {
	// 创建一个新的共享通道
	sc := NewSharedChannel(func(v int) string {
		return strconv.Itoa(v)
	})

	//监听消息, 相当于启动了多个 goroutine 来处理消息
	_ = sc.Pull(context.Background(), func(v int) bool {
		fmt.Println(v)
		return true
	})

	//推送消息
	sc.Push(1)
}

func ExampleNewSafeHashSet() {
	// 创建一个新的安全集合
	shs := NewSafeHashSet[int]()

	// 添加元素
	shs.Add(1)

	// 删除元素
	shs.Remove(1)

	// 获取元素
	shs.Contains(1)
}

func ExampleNewSharedSafeMap() {
	// 创建一个新的安全映射
	sm := NewSharedSafeMap[string]()

	// 添加元素
	sm.Store("1", "hello")

	// 获取元素
	if v, ok := sm.Load("1"); ok {
		fmt.Println(v)
	}

	// 删除元素
	sm.Del("1")
}
