package alarm_benchmark

import (
	"testing"
	"time"
	"fmt"
	"math/rand"
	"github.com/garyburd/redigo/redis"
	"strings"
)


type alarm struct {
	id            int
	expression_id int
	create_time   time.Time
	note          string
	host          string
	node          string
	filter        bool
}

//测试大量写报警的性能,报警使用hash存储
func BenchmarkWriteHashAlarms(b *testing.B) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		b.Fatal(err)
	}

	alarmCnt := 100000
	alarms := newAlarms(alarmCnt)
	size := 20

	b.ResetTimer()
	for i := 0; i < alarmCnt; i += size {
		for s := i; s < i + size; s++ {
			name := fmt.Sprintf("alarm%d", s)
			conn.Send("hmset", name, "id", alarms[s].id, "host", alarms[s].host, "filter",
				alarms[s].filter, "create_time", alarms[s].create_time.Unix(), "note",
				alarms[s].note, "node", alarms[s].node, "expression_id", alarms[s].expression_id)
			conn.Flush()
		}
	}
	conn.Receive()

	fmt.Println("end")
}

//测试大量写报警的性能,报警使用单个value存储
func BenchmarkWriteJsonAlarms(b *testing.B) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		b.Fatal(err)
	}

	alarmCnt := 100000
	size := 40

	content := strings.Repeat("sssssssssssssssssssssss", 10)
	for i := 0; i < alarmCnt; i += size {
		for s := i; s < i + size; s++ {
			name := fmt.Sprintf("alarm%d", s)
			conn.Send("set", name, content)
		}
		conn.Flush()
	}
	conn.Receive()

	fmt.Println("end")
}

//随机创建cnt个报警
func newAlarms(cnt int) []alarm {
	alarms := make([]alarm, cnt)
	for i := 0; i < cnt; i++ {
		alarms[i].id = i
		alarms[i].expression_id = rand.Int() % 50
		alarms[i].note = "sssssssssssssssssssssssss"
		alarms[i].host = fmt.Sprintf("host-host-xxx-wocao-%d", rand.Int31n(10))
		alarms[i].node = fmt.Sprintf("node%d", rand.Int31n(5))
	}
	return alarms
}
