package main

import (
	"fmt"
	"os"

	"github.com/chatroom/client/process"
)

//定义两个变量，一个表示ID；一个表示密码
var userId int
var userPwd string
var userName string

func main() {

	//接受用户输入

	var key int
	//判断是否继续显示菜单

	//var loop bool

	for true {
		fmt.Println("------欢迎登陆聊天室------")
		fmt.Println("------1 登陆聊天室------")
		fmt.Println("------2 注册用户------")
		fmt.Println("------3 退出系统------")
		fmt.Println("------请选择（1-3）------")

		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			//说明用户要登陆
			fmt.Println("请输入用户ID")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码")
			fmt.Scanf("%s\n", &userPwd)

			//完成登陆

			up := &process.UserProcess{}
			up.Login(userId, userPwd)

		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户ID")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户姓名")
			fmt.Scanf("%s\n", &userName)

			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)

		case 3:
			fmt.Println("退出系统")
			//loop = false
			os.Exit(0)
		default:
			fmt.Println("输入有误")

		}

	}
}
