package process

import (
	"fmt"

	"github.com/chatroom/client/model"
	"github.com/chatroom/common/message"
)

//我们在用户登录成功后，完成对CurUser初始化
var CurUser model.CurUser

//声明一个结构体，客户端要维护一个map,用于存储在线的用户
type UserMgr struct {
	onlineUsers map[int]*message.User
}

var (
	userMgr *UserMgr
)

//完成对userMgr初始化

func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*message.User, 1024),
	}
}

//处理返回的NotifyUserStatusMes

func (this *UserMgr) updateUserStatus(NotifyUserStatusMes *message.NotifyUserStatusMes) {

	//优化，查询本地map有无此新登录用户
	user, ok := this.onlineUsers[NotifyUserStatusMes.UserId]
	if !ok { //原来没有

		user = &message.User{

			UserId: NotifyUserStatusMes.UserId,
		}

	}

	user.UserStatus = NotifyUserStatusMes.Status
	this.onlineUsers[NotifyUserStatusMes.UserId] = user

	this.outputOnlineUsers()

}

//定义对此结构体的方法
//显示在线用户
func (this *UserMgr) outputOnlineUsers() {

	//测试
	fmt.Println(this.onlineUsers)

	//遍历userMgr
	for id := range this.onlineUsers {
		fmt.Printf("用户id：%d\t\n", id)

	}

}
