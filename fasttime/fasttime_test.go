/*
 * Copyright 2025 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package fasttime

import (
	"fmt"
	"testing"
	"time"
)

func TestUnixTimestamp(t *testing.T) {
	tsExpected := time.Now().Unix()
	ts := UnixTimestamp()
	if ts-tsExpected > 1 {
		t.Fatalf("unexpected UnixTimestamp; got %d; want %d", ts, tsExpected)
	}
}

func TestUnixDate(t *testing.T) {
	dateExpected := time.Now().Unix() / (24 * 3600)
	date := UnixDate()
	if date-dateExpected > 1 {
		t.Fatalf("unexpected UnixDate; got %d; want %d", date, dateExpected)
	}
}

func TestUnixHour(t *testing.T) {
	hourExpected := time.Now().Unix() / 3600
	hour := UnixHour()
	if hour-hourExpected > 1 {
		t.Fatalf("unexpected UnixHour; got %d; want %d", hour, hourExpected)
	}
}

func TestNow(t *testing.T) {
	fmt.Println(Now())
	tsExpected := time.Now().Unix()
	ts := Now().Unix()
	if ts-tsExpected < 1 {
		t.Fatalf("unexpected Now; got %d; want %d", ts, tsExpected)
	}
}
