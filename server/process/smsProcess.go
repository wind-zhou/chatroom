package process

import (
	"encoding/json"
	"fmt"

	"github.com/chatroom/common/message"
	"github.com/chatroom/server/utils"
)

type SmsProcess struct {
}

//处理消息的方法函数，服务器向clients群发

func (this *SmsProcess) SendGroupMes(mes *message.Message) (err error) {

	//遍历服务器端的onlineUsers map[int]*UserProcess,
	//将消息转发取出.
	//取出mes的内容 SmsMes
	var SmsMes message.SmsMes

	//这里还反序列化的目的是拿到里面的用户ID

	err = json.Unmarshal([]byte(mes.Date), &SmsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
	}

	date, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	for id, up := range userMgr.onlineUsers {

		if id == SmsMes.UserId {
			continue
		}

		tf := &utils.Transfer{
			Conn: up.Conn,
		}
		err = tf.WritePkg(date)
		if err != nil {
			fmt.Println("转发消息失败 err=", err)
		}

	}
	return

}
