package main

import (
	"testing"
	"strconv"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

func BenchmarkReadAlarms(b *testing.B) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		b.Fatal(err)
	}

	once := 100
	s := make([]interface{}, once * 2 + 1)
	s[0] = "int"

	for i := 0; i < 200; i += once {
		for j := i; j < once; j++ {
			s[j * 2 + 1] = strconv.Itoa(j + i)
			s[j * 2 + 2] = strconv.Itoa(j + i)
		}
		r, err := conn.Do("hmset", s...)
		fmt.Println(r, err, s)
	}
}