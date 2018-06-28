package clound

import (
	"fmt"
	"sort"
)

type MachineLevelDeploy struct {
	LevelConfig *MachineLevelConfig
	MachineMap  map[string]*Machine
}

func NewMachineLevelDeploy(level *MachineLevelConfig) *MachineLevelDeploy {
	p := &MachineLevelDeploy{}
	p.LevelConfig = level
	p.MachineMap = make(map[string]*Machine)

	return p
}

func (p *MachineLevelDeploy) AddMachine(m *Machine) {
	debugLog("MachineLevelDeploy.AddMachine %s %v", m.MachineId, p.LevelConfig)
	p.MachineMap[m.MachineId] = m
}

func (p *MachineLevelDeploy) RemoveMachine(machineId string) {
	debugLog("MachineLevelDeploy.RemoveMachine %s %v", machineId, p.LevelConfig)
	delete(p.MachineMap, machineId)
}

type MachineLevelDeployArray []*MachineLevelDeploy

func (p MachineLevelDeployArray) Len() int {
	return len(p)
}

func (p MachineLevelDeployArray) Less(i, j int) bool {
	return p[i].LevelConfig.Less(p[j].LevelConfig)
}

func (p MachineLevelDeployArray) Swap(i, j int) {
	temp := p[i]
	p[i] = p[j]
	p[j] = temp
}

type MachineDeployPool struct {
	MachineMap              map[string]*Machine
	MachineLevelDeployArray MachineLevelDeployArray
}

func NewMachineDeployPool() *MachineDeployPool {
	p := &MachineDeployPool{}
	p.MachineMap = make(map[string]*Machine)

	return p
}

func (p *MachineDeployPool) AddMachine(m *Machine) {
	debugLog("MachineDeployPool.AddMachine %s", m.MachineId)
	p.MachineMap[m.MachineId] = m

	var pool *MachineLevelDeploy
	for _, v := range p.MachineLevelDeployArray {
		if v.LevelConfig == m.LevelConfig {
			pool = v
			break
		}
	}
	if pool == nil {
		fmt.Println("MachineDeployPool.AddMachine new level", m.LevelConfig)
		pool = NewMachineLevelDeploy(m.LevelConfig)
		p.MachineLevelDeployArray = append(p.MachineLevelDeployArray, pool)
		sort.Sort(p.MachineLevelDeployArray)
		for _, v := range p.MachineLevelDeployArray {
			fmt.Println("    ", v.LevelConfig)
		}
	}
	pool.AddMachine(m)
}

func (p *MachineDeployPool) RemoveMachine(machineId string) {
	debugLog("MachineDeployPool.RemoveMachine %s", machineId)
	m := p.MachineMap[machineId]
	delete(p.MachineMap, machineId)

	for _, v := range p.MachineLevelDeployArray {
		if v.LevelConfig == m.LevelConfig {
			v.RemoveMachine(machineId)
			return
		}
	}
}
