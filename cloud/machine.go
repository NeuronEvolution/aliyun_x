package cloud

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

type Machine struct {
	Resource
	R                  *ResourceManagement
	Rand               *rand.Rand
	MachineId          int
	LevelConfig        *MachineLevelConfig
	InstanceArray      InstanceListSortByCostEvalDesc
	InstanceArrayCount int
	appCountCollection *AppCountCollection

	cpuCost      float64
	cpuCostValid bool
}

func NewMachine(r *ResourceManagement, machineId int, levelConfig *MachineLevelConfig) *Machine {
	m := &Machine{}
	m.R = r
	m.Rand = rand.New(rand.NewSource(0))
	m.MachineId = machineId
	m.LevelConfig = levelConfig
	m.InstanceArray = make([]*Instance, MaxInstancePerMachine)
	m.appCountCollection = NewAppCountCollection()

	return m
}

func (m *Machine) ClearInstances() {
	m.InstanceArrayCount = 0
	m.appCountCollection.Clear()
	for i := 0; i < len(m.Cpu); i++ {
		m.Cpu[i] = 0
	}
	for i := 0; i < len(m.Mem); i++ {
		m.Mem[i] = 0
	}
	m.Disk = 0
	m.P = 0
	m.M = 0
	m.PM = 0

	m.cpuCostValid = false
}

func (m *Machine) AddInstance(instance *Instance) {
	//debugLog("Machine.AddInstance %d %d", m.MachineId, instance.InstanceId)

	m.InstanceArray[m.InstanceArrayCount] = instance
	m.InstanceArrayCount++
	m.appCountCollection.Add(instance.Config.AppId)
	m.allocResource(instance)

	sort.Sort(m.InstanceArray[:m.InstanceArrayCount])
	m.R.SetInstanceDeployedMachine(instance, m)
	if m.InstanceArrayCount == 1 {
		m.R.MachineFreePool.RemoveMachine(m.MachineId)
		m.R.MachineDeployPool.AddMachine(m)
	}

	m.cpuCostValid = false

	m.calcCostEval(m.LevelConfig)

	if DebugEnabled {
		//m.debugValidation()
	}
}

func (m *Machine) RemoveInstance(instanceId int) {
	//debugLog("Machine.RemoveInstance machineId=%d,instanceId=%d", m.MachineId, instanceId)
	for i, v := range m.InstanceArray[:m.InstanceArrayCount] {
		if v.InstanceId == instanceId {
			instance := m.InstanceArray[i]
			//debugLog("Machine.RemoveInstance appId=%d", instance.Config.AppId)
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

			if m.InstanceArrayCount == 0 {
				m.R.MachineDeployPool.RemoveMachine(m.MachineId)
				m.R.MachineFreePool.AddMachine(m)
			}

			m.cpuCostValid = false

			break
		}
	}

	m.calcCostEval(m.LevelConfig)

	if DebugEnabled {
		//m.debugValidation()
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

func (m *Machine) IsEmpty() bool {
	return len(m.InstanceArray) == 0
}

func (m *Machine) ConstraintCheck(instance *Instance) bool {
	//debugLog("Machine.ConstraintCheck %s %s", m.MachineId, instance.InstanceId)

	if !constraintCheckResourceLimit(m, instance) {
		//debugLog("Machine.ConstraintCheck constraintCheckResourceLimit failed")
		return false
	}

	if !constraintCheckAppInterferenceAddInstance(
		instance.Config.AppId,
		m.appCountCollection,
		m.R.AppInterferenceConfigMap) {
		//debugLog("Machine.ConstraintCheck constraintCheckAppInterferenceAddInstance failed")
		return false
	}

	return true
}

func (m *Machine) HasBadConstraint() bool {
	return !constraintCheckAppInterference(m.appCountCollection, m.R.AppInterferenceConfigMap)
}

func (m *Machine) GetCostReal() float64 {
	totalCost := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := m.Cpu[i] / m.LevelConfig.Cpu
		if r > 0.5 {
			totalCost += 1 + 10*(math.Exp(r-0.5)-1)
		} else {
			totalCost += 1
		}
	}

	return totalCost / TimeSampleCount
}

func (m *Machine) GetCost() float64 {
	if m.cpuCostValid {
		return m.cpuCost
	}
	m.cpuCostValid = true

	totalCost := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := m.Cpu[i] / m.LevelConfig.Cpu
		if r > 0.5 {
			totalCost += 1 + 10*(Exp(r-0.5)-1)
		} else {
			totalCost += 1
		}
	}

	m.cpuCost = totalCost / TimeSampleCount

	return m.cpuCost
}

func (m *Machine) GetCostWithInstance(instance *Instance) float64 {
	totalCost := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := (m.Cpu[i] + instance.Config.Cpu[i]) / m.LevelConfig.Cpu
		if r > 0.5 {
			totalCost += 1 + 10*(Exp(r-0.5)-1)
		} else {
			totalCost += 1
		}
	}

	return totalCost / TimeSampleCount
}

func (m *Machine) GetResourceCostWithInstance(inst *Instance) float64 {
	avgCpu := float64(0)
	for i, v := range m.Cpu {
		avgCpu += v + inst.Config.Cpu[i]
	}
	avgCpu = avgCpu / float64(len(m.Cpu))

	avgMem := float64(0)
	for i, v := range m.Mem {
		avgMem += v + inst.Config.Mem[i]
	}
	avgMem = avgMem / float64(len(m.Mem))

	cpu := avgCpu / m.LevelConfig.Cpu
	mem := avgMem / m.LevelConfig.Mem
	disk := float64(m.Disk+inst.Config.Disk) / float64(m.LevelConfig.Disk)
	p := float64(m.P+inst.Config.P) / float64(m.LevelConfig.P)
	mCost := float64(m.M+inst.Config.M) / float64(m.LevelConfig.M)
	pm := float64(m.PM+inst.Config.PM) / float64(m.LevelConfig.PM)

	cost := scaleCost(cpu) +
		scaleCost(mem) +
		scaleCost(disk) +
		scaleCost(p) +
		scaleCost(mCost) +
		scaleCost(pm)

	deviation := calcResourceCostDeviation(cpu, mem, disk, p, mCost, pm)

	d := deviation - m.ResourceCostDeviation

	return cost * Exp(1+d)
}

func (m *Machine) debugValidation() {
	for i := 0; i < m.InstanceArrayCount; i++ {
		if m.InstanceArray[i] == nil {
			panic(fmt.Errorf("Machine.debugValidation machineId=%d,i=%d", m.MachineId, i))
		}
	}

	m.appCountCollection.debugValidation()
}

func (m *Machine) debugLogResource() {
	if DebugEnabled {
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
		fmt.Printf("Machine.debugLogResource %d %f %f %d %d %d %d\n",
			m.MachineId, maxCpu, maxMem, m.Disk, m.P, m.M, m.PM)
	}
}

func (m *Machine) DebugPrint() {
	fmt.Printf("Machine.DebugPrint %d %v\n", m.MachineId, m.LevelConfig)
	for i := 0; i < m.appCountCollection.ListCount; i++ {
		fmt.Printf("    %v,%v\n", m.appCountCollection.List[i],
			m.R.AppResourcesConfigMap[m.appCountCollection.List[i].AppId])
	}
	m.debugLogResource()
}
