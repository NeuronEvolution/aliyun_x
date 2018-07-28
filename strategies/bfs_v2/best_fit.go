package bfs_v2

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
)

func (s *Strategy) measureWithInstance(m *cloud.Machine, instance *cloud.Instance) (d float64) {
	disk := float64(m.Disk+instance.Config.Disk) / float64(m.LevelConfig.Disk)
	if disk > 1 {
		return math.MaxFloat64
	}

	cpuMax := float64(0)
	for i, v := range m.Cpu {
		cpu := v + instance.Config.Cpu[i]
		if cpu > cpuMax {
			cpuMax = cpu
		}
	}
	cpu := cpuMax / (m.LevelConfig.Cpu * cloud.MaxCpuRatio)
	if cpu > 1 {
		return math.MaxFloat64
	}

	memMax := float64(0)
	for i, v := range m.Mem {
		mem := v + instance.Config.Mem[i]
		if mem > memMax {
			memMax = mem
		}
	}
	mem := memMax / m.LevelConfig.Mem
	if mem > 1 {
		return math.MaxFloat64
	}

	max := float64(0)
	if disk > max {
		max = disk
	}

	if cpu > max {
		max = cpu
	}

	if mem > max {
		max = mem
	}

	return (cpu + disk + mem) * max
}

func (s *Strategy) costWithInstance(m *cloud.Machine, instance *cloud.Instance, progress float64) (cost float64) {
	return m.GetDerivationWithInstance(instance)
}

func (s *Strategy) bestFitResource(instance *cloud.Instance, cpuMax float64, progress float64) *cloud.Machine {
	min := math.MaxFloat64
	var machine *cloud.Machine
	for _, m := range s.machineDeployList {
		if progress < 0 && m.LevelConfig.Cpu != cloud.HighCpu {
			continue
		}

		if cloud.InstancesContainsApp(m.InstanceArray[:m.InstanceArrayCount], instance.Config.AppId) {
			continue
		}

		if !m.ConstraintCheckResourceLimit(instance, cpuMax) {
			continue
		}

		d := s.costWithInstance(m, instance, progress)
		if d >= min {
			continue
		}

		if !m.ConstraintCheckAppInterferenceAddInstance(instance) {
			continue
		}

		//fmt.Println(d)
		min = d
		machine = m
	}

	if machine != nil {
		//fmt.Printf("BestFitStrategy.bestFitResource cost=%f,%f,machineId=%d\n",
		//	minD, instance.Config.Cpu[0], machine.MachineId)
	}

	return machine
}

func (s *Strategy) bestFitCpuCost(instance *cloud.Instance, progress float64, all bool) *cloud.Machine {
	minCpuCost := math.MaxFloat64
	var minCpuCostMachine *cloud.Machine

	for _, m := range s.machineDeployList {
		if !all && progress < 0 && m.LevelConfig.Cpu != cloud.HighCpu {
			continue
		}

		if !all && cloud.InstancesContainsApp(m.InstanceArray[:m.InstanceArrayCount], instance.Config.AppId) {
			continue
		}

		if !m.ConstraintCheckResourceLimit(instance, 1) {
			continue
		}

		cost := m.GetLinearCostWithInstance(instance)
		if cost > minCpuCost {
			continue
		}

		if !m.ConstraintCheckAppInterferenceAddInstance(instance) {
			continue
		}

		minCpuCost = cost
		minCpuCostMachine = m
	}

	if minCpuCostMachine != nil {
		//fmt.Printf("BestFitStrategy.bestFitResource cost=%f,%f,machineId=%d\n",
		//	minCost, instance.Config.Cpu[0], minCostMachine.MachineId)
	}

	return minCpuCostMachine
}
