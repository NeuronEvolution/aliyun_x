package cloud

import (
	"bytes"
	"fmt"
	"sort"
)

type MachineLevelDeploy struct {
	LevelConfig       *MachineLevelConfig
	MachineCollection *MachineCollection
}

func NewMachineLevelDeploy(level *MachineLevelConfig) *MachineLevelDeploy {
	p := &MachineLevelDeploy{}
	p.LevelConfig = level
	p.MachineCollection = NewMachineCollection()

	return p
}

func (p *MachineLevelDeploy) AddMachine(m *Machine) {
	//debugLog("MachineLevelDeploy.AddMachine %s %v", m.MachineId, p.LevelConfig)
	p.MachineCollection.Add(m)
}

func (p *MachineLevelDeploy) RemoveMachine(machineId int) {
	//debugLog("MachineLevelDeploy.RemoveMachine %s %v", machineId, p.LevelConfig)
	p.MachineCollection.Remove(machineId)
}

type MachineLevelDeployListSortByMachineLevelDesc []*MachineLevelDeploy

func (p MachineLevelDeployListSortByMachineLevelDesc) Len() int {
	return len(p)
}

func (p MachineLevelDeployListSortByMachineLevelDesc) Less(i, j int) bool {
	return !p[i].LevelConfig.Less(p[j].LevelConfig)
}

func (p MachineLevelDeployListSortByMachineLevelDesc) Swap(i, j int) {
	temp := p[i]
	p[i] = p[j]
	p[j] = temp
}

func (p MachineLevelDeployListSortByMachineLevelDesc) First() *MachineLevelDeploy {
	if len(p) == 0 {
		return nil
	}

	return p[0]
}

func (p MachineLevelDeployListSortByMachineLevelDesc) Last() *MachineLevelDeploy {
	if len(p) == 0 {
		return nil
	}

	return p[len(p)-1]
}

type MachineDeployPool struct {
	MachineMap              map[int]*Machine
	MachineLevelDeployArray MachineLevelDeployListSortByMachineLevelDesc
}

func NewMachineDeployPool() *MachineDeployPool {
	p := &MachineDeployPool{}
	p.MachineMap = make(map[int]*Machine)

	return p
}

func (p *MachineDeployPool) AddMachine(m *Machine) {
	//debugLog("MachineDeployPool.AddMachine %d", m.MachineId)
	p.MachineMap[m.MachineId] = m

	var pool *MachineLevelDeploy
	for _, v := range p.MachineLevelDeployArray {
		if v.LevelConfig == m.LevelConfig {
			pool = v
			break
		}
	}
	if pool == nil {
		//debugLog("MachineDeployPool.AddMachine new level %v", m.LevelConfig)
		pool = NewMachineLevelDeploy(m.LevelConfig)
		p.MachineLevelDeployArray = append(p.MachineLevelDeployArray, pool)
		sort.Sort(p.MachineLevelDeployArray)
	}
	pool.AddMachine(m)
}

func (p *MachineDeployPool) RemoveMachine(machineId int) *Machine {
	//debugLog("MachineDeployPool.RemoveMachine %d", machineId)
	m, has := p.MachineMap[machineId]
	if !has {
		return nil
	}
	delete(p.MachineMap, machineId)

	for _, v := range p.MachineLevelDeployArray {
		if v.LevelConfig == m.LevelConfig {
			v.RemoveMachine(machineId)
			return m
		}
	}

	return m
}

func (p *MachineDeployPool) DebugPrint(buf *bytes.Buffer) {
	buf.WriteString("MachineDeployPool.DebugPrint\n")
	instanceCount := 0
	for _, v := range p.MachineMap {
		//v.DebugPrint()
		instanceCount += v.InstanceArrayCount
	}
	for _, v := range p.MachineLevelDeployArray {
		buf.WriteString(fmt.Sprintf("    %v machineCount=%d\n",
			v.LevelConfig, v.MachineCollection.ListCount))
	}

	cpuHighCount := 0
	cpuHighMachineList := make([]*Machine, 0)
	highCpuLimit := MaxCpu
	for _, v := range p.MachineMap {
		if v.GetCost() > highCpuLimit {
			cpuHighCount++
			cpuHighMachineList = append(cpuHighMachineList, v)
		}
	}
	sort.Sort(MachineListSortByCostDesc(cpuHighMachineList))
	for _, v := range cpuHighMachineList {
		if v.GetCost() > 1.2 {
			buf.WriteString(fmt.Sprintf("    cpuCost=%f,machineId=%d\n", v.GetCost(), v.MachineId))
		}
	}
	buf.WriteString(fmt.Sprintf("total high cpu(%f) count=%d\n", highCpuLimit, len(cpuHighMachineList)))

	buf.WriteString(fmt.Sprintf("MachineDeployPool.DebugPrint machineCount=%d,instanceCount=%d\n",
		len(p.MachineMap), instanceCount))
}
