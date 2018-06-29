package cloud

import (
	"fmt"
	"sort"
)

type MachineLevelFree struct {
	LevelConfig       *MachineLevelConfig
	MachineCollection *MachineCollection
}

func NewMachineLevelFree(level *MachineLevelConfig) *MachineLevelFree {
	p := &MachineLevelFree{}
	p.LevelConfig = level
	p.MachineCollection = NewMachineCollection()

	return p
}

func (p *MachineLevelFree) PushMachine(m *Machine) {
	//debugLog("MachineLevelFree.PushMachine %s %v", m.MachineId, m.LevelConfig)
	p.MachineCollection.Add(m)
}

func (p *MachineLevelFree) PopMachine() (m *Machine) {
	//debugLog("MachineLevelFree.PopMachine MachineListCount=%d %v",
	//p.MachineCollection.ListCount, p.LevelConfig)
	return p.MachineCollection.Pop()
}

func (p *MachineLevelFree) RemoveMachine(machineId int) {
	//debugLog("MachineLevelFree.RemoveMachine %s %v", machineId, p.LevelConfig)
	p.MachineCollection.Remove(machineId)
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
	MachineMap            map[int]*Machine
	MachineLevelFreeArray MachineLevelFreeArray
}

func NewMachineFreePool() *MachineFreePool {
	p := &MachineFreePool{}
	p.MachineMap = make(map[int]*Machine)

	return p
}

func (p *MachineFreePool) AddMachine(m *Machine) {
	//debugLog("MachineFreePool.AddMachine %s", m.MachineId)
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

func (p *MachineFreePool) RemoveMachine(machineId int) *Machine {
	//debugLog("MachineFreePool.RemoveMachine %s", machineId)
	m, has := p.MachineMap[machineId]
	if !has {
		return nil
	}

	delete(p.MachineMap, machineId)

	for _, v := range p.MachineLevelFreeArray {
		if v.LevelConfig == m.LevelConfig {
			v.RemoveMachine(machineId)
			return m
		}
	}

	return m
}

func (p *MachineFreePool) PopMachine() (m *Machine) {
	//debugLog("MachineFreePool.PopMachine")
	for _, v := range p.MachineLevelFreeArray {
		if v.MachineCollection.ListCount > 0 {
			m = v.PopMachine()
			if m != nil {
				delete(p.MachineMap, m.MachineId)
				//debugLog("MachineFreePool.PopMachine success,machineId=%s remain=%d %v",
				//	m.MachineId, len(p.MachineMap), m.LevelConfig)
				return m
			}
		}
	}

	return nil
}
