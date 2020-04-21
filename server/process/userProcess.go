package process

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/chatroom/common/message"
	"github.com/chatroom/server/model"
	"github.com/chatroom/server/utils"
)

type UserProcess struct {
	Conn net.Conn

	UserId int
}

//用来处理登陆信息请求
//将接收到的消息反序列化
//取出其中的账号密码跟预存的进行比对，并返回信息
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {

	//核心代码
	//1. 先从mes中取出mes.date,将其反序列化

	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Date), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	//1.声明一个reMes用来返回服务器信息
	var resMes message.Message

	resMes.Type = message.LoginResMesType

	//2. 声明一个LoginResMes

	var LoginResMes message.LoginResMes

	//我们需要到redis数据库去完成验证.
	//1.使用model.MyUserDao 到redis去验证

	//测试
	fmt.Println("接收到的数据", loginMes)

	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	//测试
	fmt.Println("user=", user)
	if err != nil {
		if err == model.ERROR_USER_NOTEXIST {
			LoginResMes.Code = 500
			LoginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			LoginResMes.Code = 403
			LoginResMes.Error = err.Error()
		} else {
			LoginResMes.Code = 505
			LoginResMes.Error = "服务器内部错误..."
		}
	} else {

		//用户登录放入us成功，将用户uerMgr中
		//将登录成功的userid付给this

		this.UserId = loginMes.UserId

		userMgr.AddOnlineUser(this)

		//登录成功后，将自己的状态发送给其他在线的用户

		this.NotifyOthersOnlineUser(loginMes.UserId)

		for id, _ := range userMgr.onlineUsers {
			LoginResMes.UserId = append(LoginResMes.UserId, id)
		}

		LoginResMes.Code = 200
		fmt.Println("登陆成功")
	}

	//3 将LoginResMes序列化

	date, err := json.Marshal(LoginResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//4. 将date赋值给 resMes.date

	resMes.Date = string(date)
	//5.对其序列化，准备发送
	date, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//6.发送，使用writePkg函数

	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(date)
	return
}

//用来处理注册信息
//并根据返回的信息构造返回信息并传给client
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {

	//1. 先从mes中取出mes.date,将其反序列化

	var ResgisterMes message.ResgisterMes
	err = json.Unmarshal([]byte(mes.Date), &ResgisterMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	//1.声明一个reMes用来返回服务器信息
	var resMes message.Message

	resMes.Type = message.ResgisterMesType

	//2. 声明一个LoginResMes

	var ResgisterResMes message.ResgisterResMes

	//我们需要到redis数据库去完成验证.
	//1.使用model.MyUserDao 到redis去验证

	//测试
	fmt.Println("接收到的数据", ResgisterMes)

	err = model.MyUserDao.Resgister(&ResgisterMes)

	if err != nil {
		if err == model.ERROR_USER_EXIST {
			ResgisterResMes.Code = 505
			ResgisterResMes.Error = model.ERROR_USER_EXIST.Error()
		} else {
			ResgisterResMes.Code = 506
			ResgisterResMes.Error = "注册发生未知错误..."
		}
	} else {
		ResgisterResMes.Code = 200
		fmt.Println("注册成功")
	}

	//3 将ResgisterResMes序列化

	date, err := json.Marshal(ResgisterResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//4. 将date赋值给 resMes.date

	resMes.Date = string(date)
	//5.对其序列化，准备发送
	date, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//6.发送，使用writePkg函数

	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(date)
	return

}

//编写通知所有在线的用户函数
//用户userId通知其他人自己上线
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {

	//遍历onlinesusers，然后一个个发送其他人NotifyUserStatusMes

	for id, up := range userMgr.onlineUsers {

		//过滤自己
		if id == userId {
			continue
		}

		//通知其他人

		up.NotifyMeOnline(userId)

	}
}

func (this *UserProcess) NotifyMeOnline(userId int) {

	//组装NotifyUserStatusMes

	//1.声明一个mes用来发送状态信息
	var mes message.Message

	mes.Type = message.NotifyUserStatusMesType

	//2. 声明一个NotifyUserStatusMes

	var notifyUserStatusMes message.NotifyUserStatusMes

	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//3.将其序列化
	date, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//4.将序列化后的数据赋值给mes.date

	mes.Date = string(date)

	//5.将mes反序列化

	date, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//6.发送给client

	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(date)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}

}
