package process

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/chatroom/client/utils"
	"github.com/chatroom/common/message"
)

func showMenue() {

	fmt.Println("------恭喜XX登陆------")
	fmt.Println("------1、显示用户列表------")
	fmt.Println("------2、发送信息-----")
	fmt.Println("------3、信息列表------")
	fmt.Println("------4、退出系统------")
	fmt.Println("请选择（1-4）")
	var key int
	fmt.Scanf("%d\n", &key)

	var content string
	//因为，我们总会使用到SmsProcess实例，因此我们将其定义在swtich外部
	smsProcess := &SmsProcess{}

	switch key {
	case 1:
		//fmt.Println("显示用户列表")
		userMgr.outputOnlineUsers()
	case 2:
		fmt.Println("你想对大家说的什么:")
		//用bufio做接受输入，可以接受带空格的输入
		inputReader := bufio.NewReader(os.Stdin)
		content, _ = inputReader.ReadString('\n')

		//测试
		fmt.Println(content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你即将退出退出系统")
		os.Exit(0)
	default:
		fmt.Println("你输入有误")
	}

}

//和服务器报纸通信
func serverProcessMes(conn net.Conn) {

	tf := &utils.Transfer{
		Conn: conn,
	}

	for {

		fmt.Println("客户端等待服务器发送数据")
		mes, err := tf.ReadPkg()
		if err != nil {

			fmt.Println(" readPkg err=", err)
			return
		}

		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人上线

			//取出NotifyUserStatusMes
			var NotifyUserStatusMes message.NotifyUserStatusMes

			err = json.Unmarshal([]byte(mes.Date), &NotifyUserStatusMes)

			userMgr.updateUserStatus(&NotifyUserStatusMes)
		case message.SmsMesType:
			outputGroupMes(&mes)
		default:
			fmt.Println("服务器返回未知类型")
		}

		//fmt.Println(" mes=", mes)

	}

}
