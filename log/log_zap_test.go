/*
 * Copyright 2022 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package log

import (
	"fmt"
	"testing"
)

func Test_format(t *testing.T) {
	type args struct {
		fmtArgs []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "no args",
			args: args{fmtArgs: []interface{}{"sprint"}},
			want: "sprint",
		},
		{
			name: "with 1 args",
			args: args{fmtArgs: []interface{}{"sprint", "arg1"}},
			want: "sprintarg1",
		},
		{
			name: "with 2 string args",
			args: args{fmtArgs: []interface{}{"sprint", "arg1", "arg2"}},
			want: "sprintarg1arg2",
		},
		{
			name: "with 2 string args and int args",
			args: args{fmtArgs: []interface{}{"sprint", "arg1", "arg2", 1}},
			want: "sprintarg1arg21",
		},
		{
			name: "error",
			args: args{fmtArgs: []interface{}{fmt.Errorf("error message")}},
			want: "error message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sprint(tt.args.fmtArgs...); got != tt.want {
				t.Errorf("sprint() = %v, want %v", got, tt.want)
			}
		})
	}
}
