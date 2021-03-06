package alarm_benchmark

import (
	"testing"
	"strconv"
	"strings"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

//测试大量监控配置数据在redis中的读写效率
func BenchmarkReadMonitorConfig(b *testing.B) {
	bb := strings.Repeat("ssssssssss", 5000)

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		b.Fatal(err)
	}

	size := 100
	s := make([]interface{}, size * 2 + 1)
	s[0] = "int"

	for i := 0; i < 10000; i += size {
		for j := 0; j < size; j++ {
			s[j * 2 + 1] = strconv.Itoa(j + i)
			s[j * 2 + 2] = bb
		}
		conn.Do("hmset", s...)
	}

	fmt.Println("end", err)
}
