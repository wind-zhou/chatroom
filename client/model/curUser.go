package model

import (
	"net"

	"github.com/chatroom/common/message"
)

//定义一个结构体，用来维护当前client与server的连接
//在登陆成功后对其初始化
type CurUser struct {
	Conn net.Conn
	message.User
}
