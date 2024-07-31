/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import "sync"

type HashSet[T comparable] map[T]struct{}

// NewHashSet returns an empty HashSet
func NewHashSet[T comparable]() HashSet[T] {
	return make(map[T]struct{})
}

// NewHashSetWithSize returns an empty HashSet initialized with specific size
func NewHashSetWithSize[T comparable](size int) HashSet[T] {
	return make(map[T]struct{}, size)
}

// Add adds the specified element to this set
// Always returns true due to the build-in map doesn't indicate caller whether the given element already exists
// Reserves the return type for future extension
func (s HashSet[T]) Add(value T) bool {
	s[value] = struct{}{}
	return true
}

// Contains returns true if this set contains the specified element
func (s HashSet[T]) Contains(value T) bool {
	if _, ok := s[value]; ok {
		return true
	}
	return false
}

// Remove removes the specified element from this set
// Always returns true due to the build-in map doesn't indicate caller whether the given element already exists
// Reserves the return type for future extension
func (s HashSet[T]) Remove(value T) bool {
	delete(s, value)
	return true
}

// Range calls f sequentially for each value present in the hashset.
// If f returns false, range stops the iteration.
func (s HashSet[T]) Range(f func(value T) bool) {
	for k := range s {
		if !f(k) {
			break
		}
	}
}

// Len returns the number of elements of this set
func (s HashSet[T]) Len() int {
	return len(s)
}

// SafeHashSet is a thread-safe HashSet
type SafeHashSet[T comparable] struct {
	m map[T]struct{}
	l sync.RWMutex
}

// NewSafeHashSet returns an empty HashSet
func NewSafeHashSet[T comparable]() *SafeHashSet[T] {
	return &SafeHashSet[T]{
		m: make(map[T]struct{}),
	}
}

// NewSafeHashSetWithSize returns an empty HashSet initialized with specific size
func NewSafeHashSetWithSize[T comparable](size int) *SafeHashSet[T] {
	return &SafeHashSet[T]{
		m: make(map[T]struct{}, size),
	}
}

// Add adds the specified element to this set
// Always returns true due to the build-in map doesn't indicate caller whether the given element already exists
// Reserves the return type for future extension
func (s *SafeHashSet[T]) Add(value T) bool {
	s.l.Lock()
	defer s.l.Unlock()

	s.m[value] = struct{}{}
	return true
}

// Contains returns true if this set contains the specified element
func (s *SafeHashSet[T]) Contains(value T) bool {
	s.l.RLock()
	defer s.l.RUnlock()

	if _, ok := s.m[value]; ok {
		return true
	}
	return false
}

// Remove removes the specified element from this set
// Always returns true due to the build-in map doesn't indicate caller whether the given element already exists
// Reserves the return type for future extension
func (s *SafeHashSet[T]) Remove(value T) bool {
	s.l.Lock()
	defer s.l.Unlock()

	delete(s.m, value)
	return true
}

// Range calls f sequentially for each value present in the hashset.
// If f returns false, range stops the iteration.
func (s *SafeHashSet[T]) Range(f func(value T) bool) {
	s.l.RLock()
	defer s.l.RUnlock()

	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

// Len returns the number of elements of this set
func (s *SafeHashSet[T]) Len() int {
	s.l.RLock()
	defer s.l.RUnlock()

	return len(s.m)
}
