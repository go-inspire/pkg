/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package ringbuffer

import (
	"testing"
)

func TestNewRingBuffer(t *testing.T) {
	t.Run("Create valid buffer", func(t *testing.T) {
		rb, err := NewRingBuffer(5)
		if err != nil {
			t.Fatalf("Failed to create buffer: %v", err)
		}
		if rb.Cap() != 5 {
			t.Errorf("Expected capacity 5, got %d", rb.Cap())
		}
	})

	t.Run("Create invalid buffer", func(t *testing.T) {
		_, err := NewRingBuffer(0)
		if err != ErrInvalidSize {
			t.Errorf("Expected error: %v, got: %v", ErrInvalidSize, err)
		}
	})
}

func TestPushAndShift(t *testing.T) {
	rb, _ := NewRingBuffer(3)

	t.Run("Basic push and shift", func(t *testing.T) {
		rb.Push(1)
		rb.Push(2)
		rb.Push(3)

		if val, err := rb.Shift(); val != 1 || err != nil {
			t.Errorf("Expected: 1, nil, got: %v, %v", val, err)
		}
		if val, err := rb.Shift(); val != 2 || err != nil {
			t.Errorf("Expected: 2, nil, got: %v, %v", val, err)
		}
		if val, err := rb.Shift(); val != 3 || err != nil {
			t.Errorf("Expected: 3, nil, got: %v, %v", val, err)
		}
	})

	t.Run("Shift from empty buffer", func(t *testing.T) {
		_, err := rb.Shift()
		if err != ErrBufferEmpty {
			t.Errorf("Expected error: %v, got: %v", ErrBufferEmpty, err)
		}
	})

	t.Run("Overwrite when buffer is full", func(t *testing.T) {
		rb.Push(1)
		rb.Push(2)
		rb.Push(3)
		rb.Push(4) // Overwrites 1

		val, _ := rb.Shift()
		if val != 2 {
			t.Errorf("Expected: 2, got: %v", val)
		}
	})
}

func TestFetch(t *testing.T) {
	rb, _ := NewRingBuffer(3)

	t.Run("Fetch from empty buffer", func(t *testing.T) {
		if len(rb.Fetch()) != 0 {
			t.Error("Empty buffer should return empty slice")
		}
	})

	t.Run("Fetch from non-empty buffer", func(t *testing.T) {
		rb.Push(1)
		rb.Push(2)
		rb.Push(3)

		result := rb.Fetch()
		if len(result) != 3 {
			t.Fatalf("Expected length: 3, got: %d", len(result))
		}
		if result[0] != 1 || result[1] != 2 || result[2] != 3 {
			t.Errorf("Expected: [1 2 3], got: %v", result)
		}
	})
}

func TestClear(t *testing.T) {
	rb, _ := NewRingBuffer(3)
	rb.Push(1)
	rb.Push(2)

	t.Run("Clear functionality", func(t *testing.T) {
		result := rb.Clear()
		if len(result) != 2 {
			t.Errorf("Expected length: 2, got: %d", len(result))
		}
		if !rb.IsEmpty() {
			t.Error("Buffer should be empty after clear")
		}
	})
}

func TestLenAndCap(t *testing.T) {
	rb, _ := NewRingBuffer(3)

	t.Run("Initial state", func(t *testing.T) {
		if rb.Len() != 0 || rb.Cap() != 3 {
			t.Errorf("Expected: Len=0, Cap=3, got: Len=%d, Cap=%d", rb.Len(), rb.Cap())
		}
	})

	t.Run("After adding elements", func(t *testing.T) {
		rb.Push(1)
		rb.Push(2)
		if rb.Len() != 2 {
			t.Errorf("Expected length: 2, got: %d", rb.Len())
		}
	})
}

func TestIsFullAndIsEmpty(t *testing.T) {
	rb, _ := NewRingBuffer(2)

	t.Run("Initial state", func(t *testing.T) {
		if !rb.IsEmpty() || rb.IsFull() {
			t.Error("Initial state should be empty and not full")
		}
	})

	t.Run("After adding elements", func(t *testing.T) {
		rb.Push(1)
		if rb.IsEmpty() || rb.IsFull() {
			t.Error("After adding one element, should be not empty and not full")
		}

		rb.Push(2)
		if !rb.IsFull() {
			t.Error("Buffer should be full when capacity reached")
		}
	})
}

func BenchmarkRingBuffer(b *testing.B) {
	rb, _ := NewRingBuffer(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Push(i)
		if i%2 == 0 {
			rb.Shift()
		}
	}
}
