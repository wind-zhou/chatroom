package model

//用来和数据库通信

import (
	"encoding/json"
	"fmt"

	"github.com/chatroom/common/message"
	"github.com/garyburd/redigo/redis"
)

//定义一个userDao结构体
//完成对user的操作

type UserDao struct {
	pool *redis.Pool
}

//启动服务器后初始化一个userDao实例
//做成全局变量

var (
	MyUserDao *UserDao
)

//使用构造函数，创建一个userDao实例
//花里胡哨，就是把对userDao结构体的初始化封装成了一个函数

func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{

		pool: pool,
	}
	return
}

//userDaode 方法

//1、根据用户ID返回user实例+error

func (this *UserDao) getUserId(conn redis.Conn, id int) (user *User, err error) {

	//通过id查询用户

	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		fmt.Println("redis.String err=", err)
	}

	//测试
	fmt.Println("从数据库取到的数据为：", res)

	if err != nil {
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXIST
		}
		return
	}

	user = &User{}

	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

//完成登陆校验
//如果用户密码正确，则返回user实例
//如果错误，则返回错误

func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {

	//和数据库建立链接
	//从数据库取数据
	//进行比对，并返回错误

	//1.先从链接池中哪一个链接

	conn := this.pool.Get()
	defer conn.Close()

	//2.从数据库取数据

	user, err = this.getUserId(conn, userId)
	if err != nil {
		return
	}

	//到这里确认有用户存在

	//3.进行比对，并返回错误
	if user.UserPwd != userPwd {

		err = ERROR_USER_PWD
		return

	}

	return

}

//完成注册
//如果用户密码正确，则返回user实例
//如果错误，则返回错误

func (this *UserDao) Resgister(user *message.ResgisterMes) (err error) {

	//和数据库建立链接
	//从数据库取数据
	//进行判断，并返回错误

	//1.先从链接池中哪一个链接

	conn := this.pool.Get()
	defer conn.Close()

	//2.从数据库取数据

	_, err = this.getUserId(conn, user.UserId)

	//测试
	fmt.Println("err from redis=", err)
	//如果能够取到，则说明已有此用户
	if err == nil {
		err = ERROR_USER_EXIST
		return

	}

	date, err := json.Marshal(user)
	if err != nil {
		return
	}

	//入库

	_, err = conn.Do("HSet", "users", user.UserId, string(date))
	if err != nil {
		fmt.Println("保存用户出现错误 err=", err)
		return
	}

	return

}
