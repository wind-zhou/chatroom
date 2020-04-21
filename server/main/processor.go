package main

//这是一个总控

import (
	"fmt"
	"io"
	"net"

	"github.com/chatroom/common/message"
	"github.com/chatroom/server/process"
	"github.com/chatroom/server/utils"
)

//Processor is a struct
type Processor struct {
	Conn net.Conn
}

//ServerProcessMes根据不同接受的消息类型决定调用哪个函数.
func (this *Processor) ServerProcessMes(mes *message.Message) (err error) {

	switch mes.Type {
	case message.LoginMesType:
		//处理登陆信息
		//调用serverProcessLogin方法，因此首先创建创建一个结构体
		up := &process.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)

	case message.ResgisterMesType:
		//处理注册

		up := &process.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		//创建一个SmsProcess实例完成转发群聊消息.

		smsProcess := &process.SmsProcess{}
		smsProcess.SendGroupMes(mes)

	default:
		fmt.Println("消息类型不存在")

	}
	return
}

func (this *Processor) processor2() (err error) {

	for {

		//要用readPkg方法，先创建一个结构体

		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {

				fmt.Println("客户端退出，服务端退出")
				return err

			} else {
				fmt.Println("readPkg err=", err)

				return err
			}

		}

		err = this.ServerProcessMes(&mes)
		if err != nil {

			return err

		}

	}
}
