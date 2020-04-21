package main

//这里用来初始化连接池

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

//定义全局的pool
var pool *redis.Pool

func initPool(address string, maxIdel, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdel,     //最大空闲链接数
		MaxActive:   maxActive,   //与数据库的最大连接数
		IdleTimeout: idleTimeout, //最大空闲时间
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}
