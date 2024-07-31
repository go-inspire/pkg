/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"strings"
	"testing"
)

func TestHashString_EmptyString(t *testing.T) {
	result := hashString("")
	if result == 0 {
		t.Errorf("expected non-zero hash for empty string, got %d", result)
	}
}

func TestHashString_SameStringSameHash(t *testing.T) {
	str := "test"
	hash1 := hashString(str)
	hash2 := hashString(str)
	if hash1 != hash2 {
		t.Errorf("expected same hash for same string, got %d and %d", hash1, hash2)
	}
}

func TestHashString_DifferentStringsDifferentHashes(t *testing.T) {
	hash1 := hashString("test1")
	hash2 := hashString("test2")
	if hash1 == hash2 {
		t.Errorf("expected different hashes for different strings, got %d and %d", hash1, hash2)
	}
}

func TestHashString_LongString(t *testing.T) {
	longStr := "a" + strings.Repeat("b", 1000) + "c"
	result := hashString(longStr)
	if result == 0 {
		t.Errorf("expected non-zero hash for long string, got %d", result)
	}
}
