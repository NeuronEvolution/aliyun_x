package clound

import (
	"fmt"
	"sort"
)

type MachineLevelFree struct {
	LevelConfig      *MachineLevelConfig
	MachineList      []*Machine
	MachineListCount int
}

func NewMachineLevelFree(level *MachineLevelConfig) *MachineLevelFree {
	p := &MachineLevelFree{}
	p.LevelConfig = level
	p.MachineList = make([]*Machine, 32)
	p.MachineListCount = 0

	return p
}

func (p *MachineLevelFree) PushMachine(m *Machine) {
	debugLog("MachineLevelFree.PushMachine %s %v", m.MachineId, m.LevelConfig)
	p.MachineList = append(p.MachineList, m)
	p.MachineListCount++
}

func (p *MachineLevelFree) PopMachine() (m *Machine) {
	debugLog("MachineLevelFree.PopMachine MachineListCount=%d %v", p.MachineListCount, p.LevelConfig)
	if p.MachineListCount == 0 {
		return nil
	}

	p.MachineListCount--
	m = p.MachineList[p.MachineListCount]
	p.MachineList[p.MachineListCount] = nil

	return m
}

func (p *MachineLevelFree) RemoveMachine(machineId string) {
	debugLog("MachineLevelFree.RemoveMachine %s %v", machineId, p.LevelConfig)
	for i, v := range p.MachineList {
		if v.MachineId == machineId {
			p.MachineList[i] = nil
			if p.MachineListCount > 1 && i < p.MachineListCount-1 {
				p.MachineList[i] = p.MachineList[p.MachineListCount-1]
			}

			p.MachineListCount--
		}
	}
}

type MachineLevelFreeArray []*MachineLevelFree

func (p MachineLevelFreeArray) Len() int {
	return len(p)
}

func (p MachineLevelFreeArray) Less(i, j int) bool {
	return p[i].LevelConfig.Less(p[j].LevelConfig)
}

func (p MachineLevelFreeArray) Swap(i, j int) {
	temp := p[i]
	p[i] = p[j]
	p[j] = temp
}

type MachineFreePool struct {
	MachineMap            map[string]*Machine
	MachineLevelFreeArray MachineLevelFreeArray
}

func NewMachineFreePool() *MachineFreePool {
	p := &MachineFreePool{}
	p.MachineMap = make(map[string]*Machine)

	return p
}

func (p *MachineFreePool) AddMachine(m *Machine) {
	debugLog("MachineFreePool.AddMachine %s", m.MachineId)
	p.MachineMap[m.MachineId] = m

	var pool *MachineLevelFree
	for _, v := range p.MachineLevelFreeArray {
		if v.LevelConfig == m.LevelConfig {
			pool = v
			break
		}
	}
	if pool == nil {
		fmt.Println("MachineFreePool.AddMachine new level", m.LevelConfig)
		pool = NewMachineLevelFree(m.LevelConfig)
		p.MachineLevelFreeArray = append(p.MachineLevelFreeArray, pool)
		sort.Sort(p.MachineLevelFreeArray)
		for _, v := range p.MachineLevelFreeArray {
			fmt.Println("    ", v.LevelConfig)
		}
	}
	pool.PushMachine(m)
}

func (p *MachineFreePool) RemoveMachine(machineId string) {
	debugLog("MachineFreePool.RemoveMachine %s", machineId)
	m := p.MachineMap[machineId]
	delete(p.MachineMap, machineId)

	for _, v := range p.MachineLevelFreeArray {
		if v.LevelConfig == m.LevelConfig {
			v.RemoveMachine(machineId)
			return
		}
	}
}

func (p *MachineFreePool) PopMachine() (m *Machine) {
	debugLog("MachineFreePool.PopMachine")
	for _, v := range p.MachineLevelFreeArray {
		debugLog("MachineFreePool.PopMachine %v", v.LevelConfig)
		if v.MachineListCount > 0 {
			m = v.PopMachine()
			if m != nil {
				debugLog("MachineFreePool.PopMachine success,machineId=%s %v", m.MachineId, m.LevelConfig)
				delete(p.MachineMap, m.MachineId)
				return m
			}
		}
	}

	return nil
}
