package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net"
	"runtime"
	"time"
)

var (
	addr = flag.String(
		"addr",
		"54.254.81.45:3333",
		"Address to connect.")

	duration = flag.Duration(
		"duration",
		time.Second*10,
		"")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Printf("CPU: %d\n", runtime.NumCPU())
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	max := time.Duration(0)
	min := time.Duration(math.MaxInt64)
	now := time.Now()
	buf := make([]byte, 1024)
	for i := 0; time.Now().Sub(now) < *duration; i++ {
		//准备要发送的字符串
		msg := fmt.Sprintf("ping, %03d", i)
		start := time.Now()
		n, err := conn.Write([]byte(msg))
		if err != nil {
			println("Write Buffer Error:", err.Error())
			break
		}
		//fmt.Println(msg)
		//从服务器端收字符串
		n, err = conn.Read(buf)
		if err != nil {
			println("Read Buffer Error:", err.Error())
			break
		}
		//fmt.Println(string(buf[0:n]))
		d := time.Now().Sub(start)
		if max < d {
			max = d
		}
		if min > d {
			min = d
		}

		fmt.Println(n, "bytes", "from", conn.RemoteAddr(), "min", min, "max", max, "time", d)
		//等一秒钟
		time.Sleep(time.Second)
	}

}
