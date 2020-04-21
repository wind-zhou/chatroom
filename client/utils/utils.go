package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	"github.com/chatroom/common/message"
)

//将公共方法关联到结构体
type Transfer struct {
	Conn net.Conn
	buf  [8096]byte
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {

	//buf := make([]byte, 8096)
	fmt.Println("读取客户机信息...")

	_, err = this.Conn.Read(this.buf[:4])
	if err != nil {

		return
	}

	//将buf[:4]转化从uint32

	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.buf[:4])

	//根据pkglen读取信息

	n, err := this.Conn.Read(this.buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}

	//将读取的内容进行反序列化

	err = json.Unmarshal(this.buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	return

}

func (this *Transfer) WritePkg(date []byte) (err error) {

	//发送数据给对方

	var pkgLen uint32
	pkgLen = uint32(len(date))

	//var buf [4]byte //先来个数组
	binary.BigEndian.PutUint32(this.buf[:4], pkgLen)

	//发送长度
	n, err := this.Conn.Write(this.buf[:4])
	if n != 4 || err != nil {

		fmt.Println("conn.Write err=", err)
		return
	}

	//发送date本身

	n, err = this.Conn.Write(date)
	if n != int(pkgLen) || err != nil {

		fmt.Println("conn.Write err=", err)
		return
	}
	return

}
