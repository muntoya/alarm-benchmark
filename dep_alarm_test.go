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
	alarmCnt := 10000
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

	//为报警建立索引
	idMap := make(map[int][]*alarm)
	for i := 0; i < alarmCnt; i++ {
		idMap[alarms[i].Expression_id] = append(idMap[alarms[i].Expression_id], &alarms[i])
	}

	hostMap := make(map[string][]*alarm)
	for i := 0; i < alarmCnt; i++ {
		hostMap[alarms[i].Host] = append(hostMap[alarms[i].Host], &alarms[i])
	}

	nodeMap := make(map[string][]*alarm)
	for i := 0; i < alarmCnt; i++ {
		nodeMap[alarms[i].Node] = append(nodeMap[alarms[i].Node], &alarms[i])
	}

	cpuCnt := 4

	//为规则建立索引
	depsMapMap := make([]map[int][]*dep, cpuCnt)
	for i := 0; i < cpuCnt; i++ {
		depsMapMap[i] = make(map[int][]*dep)
	}

	for i := 0; i < depCnt; i++ {
		depsMap := depsMapMap[deps[i].a_id % cpuCnt]
		depsMap[deps[i].a_id] = append(depsMap[deps[i].a_id], &deps[i])
	}

	//应用配置

	g := &sync.WaitGroup{}
	g.Add(cpuCnt)
	for i := 0; i < cpuCnt; i++ {
		//depsSlice := deps[depCnt/cpuCnt*i:depCnt/cpuCnt*(i+1)]
		go depAlarm(idMap, hostMap, nodeMap, depsMapMap[i], g)
	}
	g.Wait()

	fmt.Println("end")
}

func depAlarm(idMap map[int][]*alarm, hostMap map[string][]*alarm,
	nodeMap map[string][]*alarm, depsMap map[int][]*dep, g *sync.WaitGroup) {
	fmt.Println("len", len(depsMap))
	for aid, depList := range depsMap {
		aAlarms, ok := idMap[aid]
		if !ok {
			continue
		}
		for _, dep := range depList {
			bAlarms := idMap[dep.b_id]
			if len(bAlarms) == 0 {
				continue
			}

			switch dep.dep_type {
			case dep_host:
				for a := range aAlarms {
					for b := range bAlarms {
						if bAlarms[b].Host == aAlarms[a].Host {
							bAlarms[b].Filter = true
						}
					}
				}
			case dep_node:
				for a := range aAlarms {
					for b := range bAlarms {
						if bAlarms[b].Node == aAlarms[a].Node {
							bAlarms[b].Filter = true
						}
					}
				}
			case dep_all:
				for b := range bAlarms {
					bAlarms[b].Filter = true
				}
			default:

			}

		}
	}

	g.Done()
}
