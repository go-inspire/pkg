/*
 * Copyright 2025 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package fasttime

import (
	"sync/atomic"
	"time"
)

func init() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for tm := range ticker.C {
			t := tm.Unix()
			atomic.StoreInt64(&currentTimestamp, t)
		}
	}()
}

var currentTimestamp = time.Now().Unix()

// UnixTimestamp returns the current unix timestamp in seconds.
//
// It is faster than time.Now().Unix()
func UnixTimestamp() int64 {
	return atomic.LoadInt64(&currentTimestamp)
}

// UnixDate returns date from the current unix timestamp.
//
// The date is calculated by dividing unix timestamp by (24*3600)
func UnixDate() int64 {
	return UnixTimestamp() / (24 * 3600)
}

// UnixHour returns hour from the current unix timestamp.
//
// The hour is calculated by dividing unix timestamp by 3600
func UnixHour() int64 {
	return UnixTimestamp() / 3600
}

// Now returns the current time.
func Now() time.Time {
	return time.Unix(UnixTimestamp(), 0)
}
