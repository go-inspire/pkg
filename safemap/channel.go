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

// SharedKey is a function that generates a key for a given value.
type SharedKey[T any] func(v T) string

// SharedChannel 共享通道
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

// Push 推送到通道
func (c *SharedChannel[T]) Push(value T) {
	key := c.key(value)
	i := share(key, len(c.channels))
	c.channels[i] <- value
}

// Pull 从通道拉取数据
func (c *SharedChannel[T]) Pull(f func(value T) bool) error {
	eg, ctx := errgroup.WithContext(context.Background())
	pullFn := func(channel chan T) error {
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

func (c *SharedChannel[T]) PullContext(ctx context.Context, f func(value T) bool) error {
	eg, ctx := errgroup.WithContext(ctx)
	pullFn := func(channel chan T) error {
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
