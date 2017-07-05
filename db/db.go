package db

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

/******
type Pool struct {
    // 用来创建redis连接的方法
    Dial func() (Conn, error)
    // 如果设置了给func,那么每次p.Get()的时候都会调用该方法来验证连接的可用性
    TestOnBorrow func(c Conn, t time.Time) error
    // 定义连接池中最大连接数（超过这个数会关闭老的链接，总会保持这个数）
    MaxIdle int
    // 当前连接池中可用的链接数
    MaxActive int
    // 定义链接的超时时间，每次p.Get()的时候会检测这个连接是否超时（超时会关闭，并释放可用连接数）
    IdleTimeout time.Duration
    // 当可用连接数为0是，那么当wait=true,那么当调用p.Get()时，会阻塞等待，否则，返回nil.
    Wait bool
}
******/

var pool *redis.Pool = NewDbPool()

func NewDbPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     8,
		MaxActive:   32, // max number of connections
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "192.168.78.130:6379")
			if err != nil {
				panic(err.Error())
			}
			redis.DialConnectTimeout(time.Duration(3 * time.Second))
			redis.DialReadTimeout(time.Duration(3 * time.Second))
			redis.DialWriteTimeout(time.Duration(3 * time.Second))
			//			if _, err := c.Do("AUTH", "123456"); err != nil {
			//				c.Close()
			//				return nil, err
			//			}
			return c, err
		},
	}
}

func Exec(commandName string, args ...interface{}) (interface{}, error) {
	conn := pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}
