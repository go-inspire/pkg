/*
 * Copyright 2023 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
		res, _ := json.Marshal(r.Header)
		w.Write(res)
	})

	fmt.Println("server listening on: 80")
	if err := http.ListenAndServe(":80", nil); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

}
