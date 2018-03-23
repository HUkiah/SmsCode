package RedisContext

import (
	"github.com/garyburd/redigo/redis"
)

type RedisContext struct {
	conn redis.Conn
}

func (self *RedisContext) Open(addr string) {
	self.conn, _ = redis.Dial("tcp", addr)
}

func (self *RedisContext) Close() {
	self.conn.Close()
}

func (self *RedisContext) IsExist(number string) bool {
	is_key_exit, _ := redis.Bool(self.conn.Do("EXISTS", number))

	return is_key_exit
}

func (self *RedisContext) Read(number string) (bool, string) {

	code, err := redis.String(self.conn.Do("GET", number))
	if err != nil {
		return false, ""
	}
	return true, code

}

func (self *RedisContext) Write(number, code, time string) bool {
	_, err := self.conn.Do("SET", number, code, "EX", time)
	if err != nil {
		return false
	}
	return true
}
