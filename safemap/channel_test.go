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
	"sync"
	"testing"
)

func TestSharedChannel(t *testing.T) {
	keyFun := func(v int) string {
		return strconv.Itoa(v)
	}
	// 创建一个新的共享通道
	sc := NewSharedChannelWithSize(keyFun, 1)
	defer sc.Close()

	//监听消息, 相当于启动了多个 goroutine 来处理消息
	ctx, cancel := context.WithCancel(context.Background())

	nums := NewSafeHashSet[int]()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := sc.Pull(ctx, func(v int) bool {
			fmt.Println(v)
			if !nums.Contains(v) {
				t.Errorf("error: %d", v)
			}
			return true
		})
		fmt.Println(err)
	}()

	//推送消息
	for i := 0; i < 100; i++ {
		nums.Add(i)
		sc.Push(i)
	}
	cancel()
	wg.Wait()
}
