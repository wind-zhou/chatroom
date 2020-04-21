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
	Buf  [8096]byte
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {

	//buf := make([]byte, 8096)
	fmt.Println("读取客户机信息...")

	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {

		//fmt.Println("conn.Read err=", err)
		return
	}
	//fmt.Println("读到的buf=", buf[:4])

	//将buf[:4]转化从uint32

	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])

	//根据pkglen读取信息

	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}

	//将读取的内容进行反序列化

	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	fmt.Println("mes=", mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	return

}

func (this *Transfer) WritePkg(data []byte) (err error) {

	//先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	// 发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}

	//发送data本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	return
}
