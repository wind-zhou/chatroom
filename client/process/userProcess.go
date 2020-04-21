package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/chatroom/client/utils"
	"github.com/chatroom/common/message"
)

type UserProcess struct {
}

//写个函数，完成注册

func (this *UserProcess) Register(userId int, userPwd string, userName string) (err error) {

	//1. 连接到服务器

	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil { //拨号失败返回
		fmt.Println("net.Dial err=", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2. 准备数发送据
	var mes message.Message
	mes.Type = message.ResgisterMesType

	//3. 创建一个ResgisterMes的结构体

	var ResgisterMes message.ResgisterMes
	ResgisterMes.UserId = userId
	ResgisterMes.UserPwd = userPwd
	ResgisterMes.UserName = userName

	//4. 将ResgisterMes序列化

	date, err := json.Marshal(ResgisterMes)
	if err != nil {

		fmt.Println("json.Marshal err=", err)
		return

	}
	//5.将date付给mes.date

	mes.Date = string(date)

	//6. 将 mes序列化
	date, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//7.1 将date长度发送给服务器
	//先渠道date长度，然后转化为byte切片

	var pkgLen uint32
	pkgLen = uint32(len(date))
	var buf [4]byte //先来个数组
	binary.BigEndian.PutUint32(buf[:4], pkgLen)

	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {

		fmt.Println("conn.Write err=", err)
		return
	}
	//测试
	fmt.Printf("客户端，发送消息(注册请求)长度为=%d 内容=%s", len(date), string(date))

	//发送消息
	_, err = conn.Write(date)
	if err != nil {

		fmt.Println("conn.Write(date) err=", err)
		return
	}

	//8. 对服务器端返回的数据接收

	tf := &utils.Transfer{
		Conn: conn,
	}

	mes, err = tf.ReadPkg()
	if err != nil {

		fmt.Println(" readPkg err=", err)
		return
	}

	//9. 将mes反序列化

	var ResgisterResMes message.ResgisterResMes

	err = json.Unmarshal([]byte(mes.Date), &ResgisterResMes)
	if err != nil {

		fmt.Println(" json.Unmarshal err=", err)
		return
	}

	if ResgisterResMes.Code == 200 {

		fmt.Println("注册成功")
		fmt.Println("请重新登陆")
		os.Exit(0)
	} else if ResgisterResMes.Code == 505 {

		fmt.Println(ResgisterResMes.Error)
		return

	} else if ResgisterResMes.Code == 506 {
		fmt.Println(ResgisterResMes.Error)
		return

	}
	return

}

//写个函数，完成登陆

func (this *UserProcess) Login(userId int, userPwd string) (err error) {

	//开始定协议

	// fmt.Printf("userId=%d\n uesrPwd=%v", userId, userPw

	// return nil

	//1. 连接到服务器

	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil { //拨号失败返回
		fmt.Println("net.Dial err=", err)
		return
	}

	//延时关闭
	defer conn.Close()

	//2. 准备数发送据
	var mes message.Message
	mes.Type = message.LoginMesType

	//3. 创建一个loginMes的结构体

	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4. 将loginMes序列化

	date, err := json.Marshal(loginMes)
	if err != nil {

		fmt.Println("json.Marshal err=", err)
		return

	}

	//5.将date付给mes.date

	mes.Date = string(date)

	//6. 将 mes序列化
	date, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//7.1 将date长度发送给服务器
	//先渠道date长度，然后转化为byte切片

	var pkgLen uint32
	pkgLen = uint32(len(date))
	var buf [4]byte //先来个数组
	binary.BigEndian.PutUint32(buf[:4], pkgLen)

	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {

		fmt.Println("conn.Write err=", err)
		return
	}
	fmt.Printf("客户端，发送消息长度为=%d 内容=%s", len(date), string(date))

	//发送消息
	_, err = conn.Write(date)
	if err != nil {

		fmt.Println("conn.Write(date) err=", err)
		return
	}

	//对服务器端返回的数据接收

	tf := &utils.Transfer{
		Conn: conn,
	}

	mes, err = tf.ReadPkg()
	if err != nil {

		fmt.Println(" readPkg err=", err)
		return
	}

	//将mes反序列化

	var loginResMes message.LoginResMes

	err = json.Unmarshal([]byte(mes.Date), &loginResMes)
	if err != nil {

		fmt.Println(" json.Unmarshal err=", err)
		return
	}

	if loginResMes.Code == 200 {

		//fmt.Println("登陆成功")
		//显示当前用户列表遍历loginResMes.UserId
		for _, v := range loginResMes.UserId {
			if v == userId { //不用显示自己在线
				continue
			}
			fmt.Println("用户id:\t", v)

			//完成 客户端的 onlineUsers 完成初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			userMgr.onlineUsers[v] = user

		}

		fmt.Println()

		//初始化CurUser,用来维护链接
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline
		//开启一个协程，和服务器报纸通信

		go serverProcessMes(conn)
		for {
			showMenue()
		}
	} else if loginResMes.Code == 500 {

		fmt.Println(loginResMes.Error)

	}
	return
}
