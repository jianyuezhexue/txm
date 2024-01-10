package util

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     100,                             // 连接数量
		MaxActive:   10,                              // 最大连接数量，不确定用0
		IdleTimeout: time.Duration(60) * time.Second, // 空闲连接超时时间|查询:CONFIG GET timeout|设置CONFIG SET timeout 65
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

// 连接池中获取连接
func GetRedisConn() redis.Conn {
	return pool.Get()
}

// 将连接放回连接池
func CloseRedisConn(conn redis.Conn) {
	conn.Close()
}

// 分布式单用户锁
func Lock(conn redis.Conn, key string, time time.Duration) (bool, error) {
	return redis.Bool(conn.Do("SET", key, "1", "EX", time.Seconds(), "NX"))
}

// 释放单用户锁
func Unlock(conn redis.Conn, key string) error {
	iLoop := 3
	for {
		_, err := conn.Do("DEL", key)
		iLoop--
		if err == nil || iLoop <= 0 {
			break
		}
	}
	return nil
}
