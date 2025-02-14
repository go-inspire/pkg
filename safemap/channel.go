/*
 * Copyright 2025 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"context"
	"golang.org/x/sync/errgroup"
	"runtime"
	"sync/atomic"
)

type SharedChannelStats struct {
	*SharedStats
	PullCount uint64
}

// AddPullCount 递增
func (s *SharedChannelStats) AddPullCount() {
	atomic.AddUint64(&s.PullCount, 1)
}

func NewSharedChannelStats(n int) *SharedChannelStats {
	return &SharedChannelStats{SharedStats: NewSharedStats(n)}
}

// SharedKey 用于生成共享通道的 key, 以便将消息分组到不同的通道中.
type SharedKey[T any] func(v T) string

// SharedChannel 多 chan 的通道; 适用于多个 goroutine 同时向多个通道推送数据, 以及多个 goroutine 同时从多个通道拉取数据的场景.
// 通过对消息进行分组处理的思路, 降低了 chan 的竞争, 提高并发性能.
type SharedChannel[T any] struct {
	channels []chan T
	key      SharedKey[T]
	stats    *SharedChannelStats
}

// NewSharedChannel 创建一个SharedChannel实例，用于在多个goroutine间共享通道。
// 这个函数接受一个SharedKey参数，用于唯一标识和管理共享的通道。
// [key] 参数表示用于标识和管理共享通道的键。
// 返回值是一个指向SharedChannel实例的指针，该实例包含了多个通道和相关的统计信息。
func NewSharedChannel[T any](key SharedKey[T]) *SharedChannel[T] {
	n := runtime.GOMAXPROCS(0)
	channels := make([]chan T, n)
	for i := range channels {
		channels[i] = make(chan T)
	}
	return &SharedChannel[T]{channels: channels, key: key, stats: NewSharedChannelStats(n)}
}

// NewSharedChannelWithSize 创建一个具有指定缓冲区大小的 SharedChannel 实例。
// 这个函数根据运行时的处理器数量，为每个处理器创建一个具有相同缓冲区大小的通道。
// 参数:
//
//	key - 用于标识 SharedChannel 的键，类型为 SharedKey[T]。
//	size - 指定每个通道的缓冲区大小。
//
// 返回值:
//
//	*SharedChannel[T] - 返回一个指向新创建的 SharedChannel 实例的指针
func NewSharedChannelWithSize[T any](key SharedKey[T], size int) *SharedChannel[T] {
	n := runtime.GOMAXPROCS(0)
	channels := make([]chan T, n)
	for i := range channels {
		channels[i] = make(chan T, size)
	}
	return &SharedChannel[T]{channels: channels, key: key, stats: NewSharedChannelStats(n)}
}

// NewSharedChannelWithSharedSize 创建一个具有共享大小的 SharedChannel。
// 该函数接收一个 SharedKey，用于标识和管理共享通道中的数据类型 T。
// 参数 shared 指定有多少个共享的子通道，而 size 则指定了每个子通道的缓冲区大小。
// 返回值是一个指向 SharedChannel 结构的指针，该结构包含了共享通道的配置。
// 此函数主要用于在多个协程间共享通道，同时控制通道的大小，以实现高效的并发通信。
func NewSharedChannelWithSharedSize[T any](key SharedKey[T], shared, size int) *SharedChannel[T] {
	channels := make([]chan T, shared)
	for i := range channels {
		channels[i] = make(chan T, size)
	}
	return &SharedChannel[T]{channels: channels, key: key, stats: NewSharedChannelStats(shared)}
}

// Push 方法用于向共享通道中推送一个值。
// 该方法首先根据值计算出一个键，然后根据键选择一个通道，将值发送到该通道。
// 最后，更新统计信息，记录此次推送操作。
func (c *SharedChannel[T]) Push(value T) {
	key := c.key(value)
	i := share(key, len(c.channels))
	c.channels[i] <- value
	c.stats.AddHit(i)
}

// Pull 从通道拉取数据; 启用 goroutine 同时处理多个通道, 通过 context 控制退出. 适用于处理多个通道的场景.
// 调用该方法会 block 当前 goroutine, 直到所有通道处理完毕, 或者 context 被取消, 或者 Close 被调用.
func (c *SharedChannel[T]) Pull(ctx context.Context, f func(value T) bool) error {
	eg, ctx := errgroup.WithContext(ctx)
	pull := func(channel <-chan T) {
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return context.Cause(ctx)

				case v, ok := <-channel:
					// 通道已关闭
					if !ok {
						return nil
					}
					c.stats.AddPullCount()
					if !f(v) {
						return nil
					}
				}
			}
		})
	}
	for _, channel := range c.channels {
		pull(channel)
	}
	return eg.Wait()
}

// Close 关闭通道
func (c *SharedChannel[T]) Close() {
	for _, channel := range c.channels {
		close(channel)
	}
}

// Stats 获取通道状态
func (c *SharedChannel[T]) Stats() *SharedChannelStats {
	return c.stats
}
