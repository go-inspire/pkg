/*
 * Copyright 2022 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"runtime"
	"time"
)

var (
	addr = flag.String(
		"addr",
		"0.0.0.0:3333",
		"Address to listen.")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Printf("CPU: %d\n", runtime.NumCPU())
	flag.Parse()

	fmt.Println("listening", *addr)
	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		panic("error listening:" + err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("error Accept:" + err.Error())
		}
		fmt.Println("Accepted the connection: ", conn.RemoteAddr())

		go EchoServer(conn)
	}

}

const RECV_BUF_LEN = 1024

func EchoServer(conn net.Conn) {
	buf := make([]byte, RECV_BUF_LEN)
	defer conn.Close()

	for {
		n, err := conn.Read(buf)
		switch err {
		case nil:
			conn.Write(buf[0:n])
		case io.EOF:
			fmt.Printf("Warning: End of data: %s \n", err)
			return
		default:
			fmt.Printf("Error: Reading data : %s \n", err)
			return
		}
	}
}
