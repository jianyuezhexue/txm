package txm

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
