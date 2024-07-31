/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"testing"
)

func TestHashSet_Add(t *testing.T) {
	set := NewHashSet[int]()
	if !set.Add(1) {
		t.Errorf("expected true, got false")
	}
	if !set.Contains(1) {
		t.Errorf("expected set to contain 1")
	}
}

func TestHashSet_Contains(t *testing.T) {
	set := NewHashSet[int]()
	set.Add(1)
	if !set.Contains(1) {
		t.Errorf("expected set to contain 1")
	}
	if set.Contains(2) {
		t.Errorf("expected set to not contain 2")
	}
}

func TestHashSet_Remove(t *testing.T) {
	set := NewHashSet[int]()
	set.Add(1)
	if !set.Remove(1) {
		t.Errorf("expected true, got false")
	}
	if set.Contains(1) {
		t.Errorf("expected set to not contain 1")
	}
}

func TestHashSet_Range(t *testing.T) {
	set := NewHashSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	sum := 0
	set.Range(func(value int) bool {
		sum += value
		return true
	})

	if sum != 6 {
		t.Errorf("expected sum to be 6, got %d", sum)
	}
}

func TestHashSet_Len(t *testing.T) {
	set := NewHashSet[int]()
	if set.Len() != 0 {
		t.Errorf("expected length to be 0, got %d", set.Len())
	}
	set.Add(1)
	if set.Len() != 1 {
		t.Errorf("expected length to be 1, got %d", set.Len())
	}
	set.Add(2)
	if set.Len() != 2 {
		t.Errorf("expected length to be 2, got %d", set.Len())
	}
	set.Remove(1)
	if set.Len() != 1 {
		t.Errorf("expected length to be 1, got %d", set.Len())
	}
}
