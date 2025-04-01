/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

// Package ringbuffer 实现了简单的环形缓冲区数据结构。
package ringbuffer

import (
	"errors"
)

var (
	// ErrBufferEmpty 表示缓冲区为空时发生的错误
	ErrBufferEmpty = errors.New("ring buffer: buffer is empty")
	// ErrInvalidSize 表示创建缓冲区时指定了无效的大小
	ErrInvalidSize = errors.New("ring buffer: invalid buffer size")
)

// RingBuffer 表示一个环形缓冲区数据结构。
// 使用固定大小的切片实现，通过头尾指针管理数据。
type RingBuffer struct {
	head int           // 头指针，指向第一个有效元素
	tail int           // 尾指针，指向下一个插入位置
	len  int           // 当前缓冲区中的元素数量
	size int           // 缓冲区的总容量
	buf  []interface{} // 底层数据存储
}

// NewRingBuffer 创建并返回一个新的 RingBuffer 实例。
// 参数 size 指定缓冲区的容量大小。
// 如果 size <= 0 会返回 ErrInvalidSize 错误。
func NewRingBuffer(size int) (*RingBuffer, error) {
	if size <= 0 {
		return nil, ErrInvalidSize
	}

	rb := &RingBuffer{
		buf:  make([]interface{}, size),
		size: size,
	}
	return rb, nil
}

// Push 向缓冲区尾部添加一个元素。
// 如果缓冲区已满，最旧的元素会被覆盖。
func (rb *RingBuffer) Push(val interface{}) {
	rb.buf[rb.tail] = val
	rb.tail = (rb.tail + 1) % rb.size

	if rb.len < rb.size {
		rb.len++
	} else {
		// 缓冲区已满，移动头指针
		rb.head = rb.tail
	}
}

// Shift 从缓冲区头部移除并返回一个元素。
// 如果缓冲区为空，返回 ErrBufferEmpty 错误。
func (rb *RingBuffer) Shift() (interface{}, error) {
	if rb.len <= 0 {
		return nil, ErrBufferEmpty
	}

	val := rb.buf[rb.head]
	rb.head = (rb.head + 1) % rb.size
	rb.len--
	return val, nil
}

// Fetch 返回缓冲区中的所有元素，但不移除它们。
// 返回的元素顺序与插入顺序一致。
func (rb *RingBuffer) Fetch() []interface{} {
	if rb.len == 0 {
		return []interface{}{}
	}

	result := make([]interface{}, rb.len)
	for i := 0; i < rb.len; i++ {
		result[i] = rb.buf[(rb.head+i)%rb.size]
	}

	return result
}

// Clear 清空缓冲区并返回所有元素。
func (rb *RingBuffer) Clear() []interface{} {
	result := rb.Fetch()

	rb.head = 0
	rb.tail = 0
	rb.len = 0

	return result
}

// Len 返回当前缓冲区中的元素数量。
func (rb *RingBuffer) Len() int {
	return rb.len
}

// Cap 返回缓冲区的总容量。
func (rb *RingBuffer) Cap() int {
	return rb.size
}

// IsFull 检查缓冲区是否已满。
func (rb *RingBuffer) IsFull() bool {
	return rb.len == rb.size
}

// IsEmpty 检查缓冲区是否为空。
func (rb *RingBuffer) IsEmpty() bool {
	return rb.len == 0
}
