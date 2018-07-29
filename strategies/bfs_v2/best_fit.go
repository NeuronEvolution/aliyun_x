package bfs_v2

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
)

func (s *Strategy) skipLow(instance *cloud.Instance) bool {
	if !s.isMem8(instance) {
		return true
	}

	if instance.Config.CpuAvg >= 8 || instance.Config.CpuAvg >= 6 {
		return true
	}

	return false
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

		//if s.skipLow(instance) {
		//	continue
		//}

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

		//if !all && s.skipLow(instance) {
		//	continue
		//}

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
