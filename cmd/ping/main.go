/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	times     = flag.Int("n", 10, "ping times")
	targetUrl = flag.String("url", "https://www.okx.com/api/v5/public/time", "ping url")
)

func main() {
	flag.Parse()

	//最大,最小,平均
	var (
		max   = time.Duration(0)
		min   = time.Duration(1<<63 - 1)
		total = time.Duration(0)
	)
	n := 0
	url := *targetUrl
	//预热
	ping(url)
	for i := 0; i < *times; i++ {
		start := time.Now()

		t1 := ping(url)
		t2 := ping(url)
		//serverTimes += t2 - t1
		timeSpen := t2 - t1
		d1 := time.Since(start)
		d2 := time.Duration(timeSpen) * time.Millisecond
		diff := d1 - d2
		total += diff
		n++
		if diff > max {
			max = diff
		}
		if diff < min {
			min = diff
		}
		fmt.Println("ping", url, d1, d2, d1-d2)
		time.Sleep(time.Second)
	}
	fmt.Println("ping", url, "Max:", max, "Min:", min, "Avg:", total/time.Duration(n))
}

func ping(url string) int64 {
	//fmt.Println("ping", url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("ping error", err)
		return 0
	}
	defer res.Body.Close()
	//{"code":"0","data":[{"ts":"1726214671723"}],"msg":""}
	var body struct {
		Code string `json:"code"`
		Data []struct {
			Ts string `json:"ts"`
		}
	}
	json.NewDecoder(res.Body).Decode(&body)
	if body.Code != "0" {
		fmt.Println("ping error", body.Code)
		return 0
	}
	t, _ := strconv.ParseInt(body.Data[0].Ts, 10, 64)
	return t
}
