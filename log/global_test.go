/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"testing"
)

func TestNamed(t *testing.T) {
	cfg = Config{
		DefaultLevel: InfoLevel,
		Named: map[string]Level{
			"a.b.c": DebugLevel,
			"a.b":   InfoLevel,
			"a":     WarnLevel,
		},
	}

	tests := []struct {
		name   string
		input  string
		expect Level
	}{
		{
			name:   "a.b.c->a.b.c",
			input:  "a.b.c",
			expect: DebugLevel,
		},
		{
			name:   "a,b,c.d->a.b.c",
			input:  "a.b.c.d",
			expect: DebugLevel,
		},
		{
			name:   "a.b->a.b",
			input:  "a.b",
			expect: InfoLevel,
		},
		{
			name:   "a.b.d->a.b",
			input:  "a.b.d",
			expect: InfoLevel,
		},
		{
			name:   "a->a",
			input:  "a",
			expect: WarnLevel,
		},
		{
			name:   "b->default",
			input:  "b",
			expect: InfoLevel,
		},
		{
			name:   "a.b.d->a.b",
			input:  "a.b.d",
			expect: InfoLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Named(tt.input).Level(); got != tt.expect {
				t.Errorf("Named(%s).Level() = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}
