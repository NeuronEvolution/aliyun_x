package cloud

import (
	"fmt"
	"math"
	"sort"
)

type Machine struct {
	R                  *ResourceManagement
	MachineId          int
	LevelConfig        *MachineLevelConfig
	InstanceArray      InstanceArray
	InstanceArrayCount int
	appCountCollection *AppCountCollection
	Cpu                [TimeSampleCount]float64 //todo decimal
	Mem                [TimeSampleCount]float64 //todo decimal
	Disk               int
	P                  int
	M                  int
	PM                 int
}

func NewMachine(r *ResourceManagement, machineId int, levelConfig *MachineLevelConfig) *Machine {
	m := &Machine{}
	m.R = r
	m.MachineId = machineId
	m.LevelConfig = levelConfig
	m.appCountCollection = NewAppCountCollection()

	return m
}

func (m *Machine) AddInstance(instance *Instance) error {
	//debugLog("Machine.AddInstance %s %s", m.MachineId, instance.InstanceId)
	if !m.ConstraintCheck(instance) {
		return fmt.Errorf("%s add %s ConstraintCheck failed", m.MachineId, instance.InstanceId)
	}

	m.InstanceArray = append(m.InstanceArray, instance)
	m.InstanceArrayCount++
	m.appCountCollection.Add(instance.Config.AppId)
	m.allocResource(instance)

	sort.Sort(m.InstanceArray)

	return nil
}

func (m *Machine) RemoveInstance(instanceId int) {
	debugLog("Machine.RemoveInstance", m.MachineId, instanceId)
	for i, v := range m.InstanceArray {
		if v.InstanceId == instanceId {
			instance := m.InstanceArray[i]
			m.InstanceArray[i] = nil
			if m.InstanceArrayCount > 1 && i < m.InstanceArrayCount-1 {
				for j := i; j < m.InstanceArrayCount-1; j++ {
					m.InstanceArray[j] = m.InstanceArray[j+1]
				}
				m.InstanceArray[m.InstanceArrayCount-1] = nil
			}

			m.InstanceArrayCount--
			m.appCountCollection.Remove(instance.Config.AppId)
			m.freeResource(instance)

			break
		}
	}
}

func (m *Machine) IsEmpty() bool {
	return len(m.InstanceArray) == 0
}

func (m *Machine) ConstraintCheck(instance *Instance) bool {
	//debugLog("Machine.ConstraintCheck %s %s", m.MachineId, instance.InstanceId)

	if !constraintCheckAppInterference(
		instance.Config.AppId,
		m.appCountCollection,
		m.R.AppInterferenceConfigMap) {
		//debugLog("Machine.ConstraintCheck constraintCheckAppInterference failed")
		return false
	}

	if !constraintCheckResourceLimit(m, instance) {
		//debugLog("Machine.ConstraintCheck constraintCheckResourceLimit failed")
		return false
	}

	return true
}

func (m *Machine) debugLogResource() {
	if debugEnabled {
		maxCpu := float64(0)
		for _, v := range m.Cpu {
			if v > maxCpu {
				maxCpu = v
			}
		}
		maxMem := float64(0)
		for _, v := range m.Mem {
			if v > maxMem {
				maxMem = v
			}
		}
		debugLog("Machine.debugLogResource %s %f %f %d %d %d %d",
			m.MachineId, maxCpu, maxMem, m.Disk, m.P, m.M, m.PM)
	}
}

func (m *Machine) allocResource(instance *Instance) {
	c := instance.Config
	for i, v := range c.Cpu {
		m.Cpu[i] += v
	}
	for i, v := range c.Mem {
		m.Mem[i] += v
	}
	m.Disk += c.Disk
	m.M += c.M
	m.P += c.P
	m.PM += c.PM

	//m.debugLogResource()
}

func (m *Machine) freeResource(instance *Instance) {
	c := instance.Config
	for i, v := range c.Cpu {
		m.Cpu[i] -= v
	}
	for i, v := range c.Mem {
		m.Mem[i] -= v
	}
	m.Disk -= c.Disk
	m.M -= c.M
	m.P -= c.P
	m.PM -= c.PM

	//m.debugLogResource()
}

func (m *Machine) DebugPrint() {
	debugLog("Machine.DebugPrint %s %v", m.MachineId, m.LevelConfig)
	for i := 0; i < m.appCountCollection.ListCount; i++ {
		debugLog("    %v", m.appCountCollection.List[i])
	}
	m.debugLogResource()
}

func (m *Machine) CalculateCost() float64 {
	totalCost := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		s := 1 + 10*(math.Exp(math.Max(0, m.Cpu[i]/m.LevelConfig.Cpu-0.5))-1)
		totalCost += s
	}

	return totalCost / TimeSampleCount
}
