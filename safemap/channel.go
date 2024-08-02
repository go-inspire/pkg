/*
 * Copyright (c) 2024 JQuant Authors. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"context"
	"golang.org/x/sync/errgroup"
	"runtime"
)

// SharedKey 用于生成共享通道的 key, 以便将消息分组到不同的通道中.
type SharedKey[T any] func(v T) string

// SharedChannel 多 chan 的通道; 适用于多个 goroutine 同时向多个通道推送数据, 以及多个 goroutine 同时从多个通道拉取数据的场景.
// 通过对消息进行分组处理的思路, 降低了 chan 的竞争, 提高并发性能.
type SharedChannel[T any] struct {
	channels []chan T
	key      SharedKey[T]
}

// NewSharedChannel 创建一个新的共享通道
func NewSharedChannel[T any](key SharedKey[T]) *SharedChannel[T] {
	n := runtime.GOMAXPROCS(0)
	channels := make([]chan T, n)
	for i := range channels {
		channels[i] = make(chan T)
	}
	return &SharedChannel[T]{channels: channels, key: key}
}

// NewSharedChannelWithSize 创建一个新的共享通道
func NewSharedChannelWithSize[T any](key SharedKey[T], size int) *SharedChannel[T] {
	n := runtime.GOMAXPROCS(0)
	channels := make([]chan T, n)
	for i := range channels {
		channels[i] = make(chan T, size)
	}
	return &SharedChannel[T]{channels: channels, key: key}
}

// Push 推送消息到通道
func (c *SharedChannel[T]) Push(value T) {
	key := c.key(value)
	i := share(key, len(c.channels))
	c.channels[i] <- value
}

// Pull 从通道拉取数据; 启用 goroutine 同时处理多个通道, 通过 context 控制退出. 适用于处理多个通道的场景.
// 调用该方法会 block 当前 goroutine, 直到所有通道处理完毕, 或者 context 被取消, 或者 Close 被调用.
func (c *SharedChannel[T]) Pull(ctx context.Context, f func(value T) bool) error {
	eg, ctx := errgroup.WithContext(ctx)
	pullFn := func(channel <-chan T) error {
		for {
			select {
			case <-ctx.Done():
				return context.Cause(ctx)

			case v := <-channel:
				if !f(v) {
					return nil
				}
			}
		}
	}
	for _, channel := range c.channels {
		eg.Go(func() error {
			return pullFn(channel)
		})
	}
	return eg.Wait()
}

// Close 关闭通道
func (c *SharedChannel[T]) Close() {
	for _, channel := range c.channels {
		close(channel)
	}
}
