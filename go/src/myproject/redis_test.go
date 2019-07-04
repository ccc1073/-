package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main() {

	conn, err := redis.Dial("tcp", ":6379/2")

	if err != nil {
		println("连接错误")
		return
	}
	defer conn.Close()

	// 执行操作数据库语句
	conn.Send("set", "123456", "haha")

	conn.Flush()

	rep, _ := conn.Receive()

	fmt.Println(rep)

}
