package alarm_benchmark

import (
	"testing"
	//	"strconv"
	//	"strings"
	"fmt"
	"math/rand"
	"sync"
)

const (
	dep_host int = iota
	dep_node     = iota
	dep_all      = iota
)

type dep struct {
	a_id     int
	b_id     int
	dep_type int
}

func BenchmarkDepAlarm(b *testing.B) {
	alarmCnt := 100000
	depCnt := 512

	//初始化报警
	alarms := newAlarms(alarmCnt)

	//初始化依赖规则
	deps := make([]dep, depCnt)
	for i := 0; i < depCnt; i++ {
		deps[i].a_id = rand.Int() % 50
		deps[i].b_id = rand.Int() % 50
		deps[i].dep_type = rand.Int() % 3
	}

	b.ResetTimer()

	//建立索引
	idMap := make(map[int][]*alarm)
	for i := 0; i < alarmCnt; i++ {
		idMap[alarms[i].expression_id] = append(idMap[alarms[i].expression_id], &alarms[i])
	}

	hostMap := make(map[string][]*alarm)
	for i := 0; i < alarmCnt; i++ {
		hostMap[alarms[i].host] = append(hostMap[alarms[i].host], &alarms[i])
	}

	nodeMap := make(map[string][]*alarm)
	for i := 0; i < alarmCnt; i++ {
		nodeMap[alarms[i].node] = append(nodeMap[alarms[i].node], &alarms[i])
	}

	//应用配置
	cpuCnt := 4
	g := &sync.WaitGroup{}
	g.Add(cpuCnt)
	for i := 0; i < cpuCnt; i++ {
		depsSlice := deps[depCnt/cpuCnt*i:depCnt/cpuCnt*(i+1)]
		go depAlarm(depsSlice, alarms, idMap, hostMap, nodeMap, g)
	}
	g.Wait()

	fmt.Println("end")
}

func depAlarm(deps []dep, alarms []alarm, idMap map[int][]*alarm, hostMap map[string][]*alarm,
	nodeMap map[string][]*alarm, g *sync.WaitGroup) {
	for _, dep := range deps {
		aAlarms, ok := idMap[dep.a_id]
		if !ok {
			continue
		}

		bAlarms, ok := idMap[dep.b_id]
		if !ok {
			continue
		}

		switch dep.dep_type {
		case dep_host:
			for a := range aAlarms {
				for b := range bAlarms {
					if bAlarms[b].host == aAlarms[a].host {
						bAlarms[b].filter = true
					}
				}
			}
		case dep_node:
			for a := range aAlarms {
				for b := range bAlarms {
					if bAlarms[b].node == aAlarms[a].node {
						bAlarms[b].filter = true
					}
				}
			}
		case dep_all:
			for b := range bAlarms {
				bAlarms[b].filter = true
			}
		default:

		}
	}
	g.Done()
}
