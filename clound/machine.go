package clound

import (
	"fmt"
	"sort"
)

type Machine struct {
	ResourceManagement *ResourceManagement
	MachineId          string
	LevelConfig        *MachineLevelConfig
	InstanceArray      InstanceArray
	InstanceArrayCount int
	appCountMap        map[string]int
}

func NewMachine(r *ResourceManagement, machineId string, levelConfig *MachineLevelConfig) *Machine {
	m := &Machine{}
	m.ResourceManagement = r
	m.MachineId = machineId
	m.LevelConfig = levelConfig
	m.appCountMap = make(map[string]int)

	return m
}

func (m *Machine) AddInstance(instance *Instance) error {
	debugLog("Machine.AddInstance %s %s", m.MachineId, instance.InstanceId)
	if !m.constraintCheck(instance) {
		return fmt.Errorf("%s add %s constraintCheck failed", m.MachineId, instance.InstanceId)
	}

	m.InstanceArray = append(m.InstanceArray, instance)
	m.InstanceArrayCount++

	sort.Sort(m.InstanceArray)

	return nil
}

func (m *Machine) RemoveInstance(instanceId string) {
	debugLog("Machine.RemoveInstance", m.MachineId, instanceId)
	for i, v := range m.InstanceArray {
		if v.InstanceId == instanceId {
			m.InstanceArray[i] = nil
			if m.InstanceArrayCount > 1 && i < m.InstanceArrayCount-1 {
				for j := i; j < m.InstanceArrayCount-1; j++ {
					m.InstanceArray[j] = m.InstanceArray[j+1]
				}
				m.InstanceArray[m.InstanceArrayCount-1] = nil
			}

			m.InstanceArrayCount--

			break
		}
	}
}

func (m *Machine) IsEmpty() bool {
	return len(m.InstanceArray) == 0
}

func (m *Machine) constraintCheck(instance *Instance) bool {
	debugWrite("Machine.constraintCheck %s %s", m.MachineId, instance.InstanceId)

	if !constraintCheckAppInterference(
		instance.Config.AppId,
		m.appCountMap,
		m.ResourceManagement.appInterferenceConfigMap) {
		debugLog("Machine.ConstraintCheck constraintCheckAppInterference failed")
		return false
	}

	if !constraintCheckResourceLimit(m, InstanceArray{}) {
		debugLog("Machine.ConstraintCheck constraintCheckResourceLimit failed")
		return false
	}

	return true
}
