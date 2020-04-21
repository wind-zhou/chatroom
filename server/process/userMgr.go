package process

import "fmt"

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

var (
	userMgr *UserMgr
)

//完成对userMgr初始化

func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//接下来定义对此切片的几种操作

//完成对onlineUser添加

func (this *UserMgr) AddOnlineUser(up *UserProcess) {

	this.onlineUsers[up.UserId] = up
}

//删除
func (this *UserMgr) DelOnlineUser(userId int) {

	delete(this.onlineUsers, userId)

}

//返回当前在线用户

func (this *UserMgr) GetOnlineUser() map[int]*UserProcess {
	return this.onlineUsers

}

//根据id返回值

func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	//从map中取值
	up, ok := this.onlineUsers[userId]
	if !ok { //没有此用
		err = fmt.Errorf("用户%d不存在", userId)
		return
	}
	return
}
