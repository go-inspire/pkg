/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"sync"
)

// HashSet 通过 map 实现唯一性的集合, 非线程安全.
type HashSet[T comparable] map[T]struct{}

// NewHashSet 返回一个空的 HashSet
func NewHashSet[T comparable]() HashSet[T] {
	return make(map[T]struct{})
}

// NewHashSetWithSize 返回一个空的 HashSet, 并初始化指定大小.
func NewHashSetWithSize[T comparable](size int) HashSet[T] {
	return make(map[T]struct{}, size)
}

// Add 将指定元素添加到此集合中
// 由于内置 map 不会指示调用者给定元素是否已存在，因此总是返回 true
// 保留返回类型以供将来扩展
func (s HashSet[T]) Add(value T) bool {
	s[value] = struct{}{}
	return true
}

// Contains 如果此集合包含指定元素，则返回 true
func (s HashSet[T]) Contains(value T) bool {
	if _, ok := s[value]; ok {
		return true
	}
	return false
}

// Remove 从此集合中删除指定元素.
func (s HashSet[T]) Remove(value T) bool {
	delete(s, value)
	return true
}

// Range 为此集合中的每个值调用 f
func (s HashSet[T]) Range(f func(value T) bool) {
	for k := range s {
		if !f(k) {
			break
		}
	}
}

// Len 返回此集合的元素数量
func (s HashSet[T]) Len() int {
	return len(s)
}

// SafeHashSet 是线程安全的 HashSet
type SafeHashSet[T comparable] struct {
	dirty sync.Map
}

// NewSafeHashSet 返回一个空的 SafeHashSet
func NewSafeHashSet[T comparable]() *SafeHashSet[T] {
	return &SafeHashSet[T]{
		dirty: sync.Map{},
	}
}

// Add 将指定元素添加到此集合中
func (s *SafeHashSet[T]) Add(value T) bool {
	_, loaded := s.dirty.LoadOrStore(value, struct{}{})
	return !loaded
}

// Contains 如果此集合包含指定元素，则返回 true
func (s *SafeHashSet[T]) Contains(value T) bool {
	_, loaded := s.dirty.Load(value)
	return loaded
}

// Remove 从此集合中删除指定元素.
func (s *SafeHashSet[T]) Remove(value T) bool {
	_, loaded := s.dirty.LoadAndDelete(value)
	return loaded
}

// Range 为此集合中的每个值调用 f
func (s *SafeHashSet[T]) Range(f func(value T) bool) {
	s.dirty.Range(func(key, _ interface{}) bool {
		return f(key.(T))
	})
}

// Len 返回此集合的元素数量
func (s *SafeHashSet[T]) Len() int {
	l := 0
	s.dirty.Range(func(_, _ interface{}) bool {
		l++
		return true
	})
	return l
}
