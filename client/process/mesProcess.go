package process

import (
	"encoding/json"
	"fmt"

	"github.com/chatroom/client/utils"
	"github.com/chatroom/common/message"
)

type SmsProcess struct {
}

//发送群聊信息

func (this *SmsProcess) SendGroupMes(content string) (err error) {
	//1.定义mes结构体
	var mes message.Message
	mes.Type = message.SmsMesType

	//2.定义SmsMes结构体

	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	//3.序列化smsmes
	date, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes err=", err)
		return
	}

	//4.将date赋值给mes.date

	mes.Date = string(date)

	//5.将mes序列化

	date, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes err=", err)
		return
	}

	//6.发送给服务器

	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}

	err = tf.WritePkg(date)
	if err != nil {
		fmt.Println("SendGroupMes err=", err)
		return
	}
	return

}
