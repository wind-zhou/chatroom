package main

import (
	"fmt"
	"net"
	"time"

	"github.com/chatroom/server/model"
)

func Process(conn net.Conn) {
	//延时关闭coon
	defer conn.Close()

	//这里创建一个总控，并调用

	processor := &Processor{
		Conn: conn,
	}
	err := processor.processor2()
	if err != nil {
		fmt.Println("客户端协程发生错误")
		return
	}

}

//编写一个函数，完成对userDao初始化

func initUserDao() {

	model.MyUserDao = model.NewUserDao(pool)

}
func main() {

	//服务器启动时，初始化redis连接池

	initPool("localhost:6379", 16, 0, 300*time.Second)

	//初始化userDao
	initUserDao()

	//提示信息
	fmt.Println("----服务器监听8889----")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}

	//一旦监听成功，就等待客户端来连接

	for {

		fmt.Println("----等待客户段连接----")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}

		//若连接成功，则起一个协程
		fmt.Println("连接成功")

		go Process(conn)
	}

}
