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
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = sc.Pull(ctx, func(v int) bool {
			fmt.Println(v)
			return true
		})
	}()

	//推送消息
	sc.Push(1)
	cancel()
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

	// 获取元素个数
	shs.Len()

	// 遍历元素
	shs.Range(func(v int) bool {
		fmt.Println(v)
		return true
	})
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
	sm.Delete("1")

	// 获取元素个数
	sm.Len()

	// 遍历元素
	sm.Range(func(k string, v string) bool {
		fmt.Println(k, v)
		return true
	})
}
