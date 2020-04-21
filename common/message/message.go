package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	ResgisterMesType        = "ResgisterMes"
	ResgisterResMesType     = "ResgisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

//这里我们定义几个用户状态的常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` //消息类型
	Date string `json:"date"` //数据内容
}

//发送时的消息类型
type LoginMes struct {
	UserId   int    `json:"userid"`   //用户ID
	UserPwd  string `json:"userpwd"`  //用户密码
	UserName string `json:"username"` //用户名
}

//返回的消息类型
type LoginResMes struct {
	Code   int    `json:"code"`   //返回状态码
	Error  string `json:"error"`  //错误信息
	UserId []int  `json:"userid"` //在线用户id
}

//注册发送时的消息类型
type ResgisterMes struct {
	UserId   int    `json:"userid"`   //用户ID
	UserPwd  string `json:"userpwd"`  //用户密码
	UserName string `json:"username"` //用户名
}

//注册返回的类型
type ResgisterResMes struct {
	Code  int    `json:"code"`  //返回状态码
	Error string `json:"error"` //错误信息
}

//为了配合服务器端推送用户上线变化
type NotifyUserStatusMes struct {
	UserId int `json:"userid"`
	Status int `json:"status"`
}

//定义一个通信时发送消息类型

type SmsMes struct {
	Content string `json;"content"`
	User           //匿名结构体

}
