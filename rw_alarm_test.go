package alarm_benchmark

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type alarm struct {
	Id            int       `json:"id"`
	Expression_id int       `json:"expression_id"`
	Create_time   time.Time `json:"create_time"`
	Note          string    `json:"note"`
	Host          string    `json:"host"`
	Node          string    `json:"node"`
	Filter        bool      `json:"filter"`
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
		for s := i; s < i+size; s++ {
			name := fmt.Sprintf("alarm%d", s)
			conn.Send("hmset", name, "id", alarms[s].Id, "host", alarms[s].Host, "filter",
				alarms[s].Filter, "create_time", alarms[s].Create_time.Unix(), "note",
				alarms[s].Note, "node", alarms[s].Node, "expression_id", alarms[s].Expression_id)
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
	alarms := newAlarms(alarmCnt)
	size := 40

	for i := 0; i < alarmCnt; i += size {
		for s := i; s < i+size; s++ {
			name := fmt.Sprintf("alarm%d", s)
			content, _ := json.Marshal(&alarms[s])
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
		alarms[i].Id = i
		alarms[i].Expression_id = rand.Int() % 50
		alarms[i].Note = "sssssssssssssssssssssssss"
		alarms[i].Host = fmt.Sprintf("host-host-xxx-wocao-%d", rand.Int31n(10))
		alarms[i].Node = fmt.Sprintf("node%d", rand.Int31n(5))
	}
	return alarms
}
