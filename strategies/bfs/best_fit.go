package bfs

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
)

func (s *BestFitStrategy) bestFitResource(
	instance *cloud.Instance, skip *cloud.Machine, cpuMax float64) *cloud.Machine {

	minResourceCost := math.MaxFloat64
	var minResourceCostMachine *cloud.Machine

	for _, m := range s.machineDeployList {
		if skip != nil && m.MachineId == skip.MachineId {
			continue
		}

		resourceCost := m.GetResourceCostWithInstance(instance)
		if resourceCost >= minResourceCost {
			continue
		}

		cost := m.GetCostWithInstance(instance)
		if cost > cpuMax {
			continue
		}

		if m.ConstraintCheck(instance) {
			minResourceCost = resourceCost
			minResourceCostMachine = m
		}
	}

	if minResourceCostMachine != nil {
		//fmt.Printf("BestFitStrategy.bestFitResource cost=%f,%f,machineId=%d\n",
		//	minCost, instance.Config.Cpu[0], minCostMachine.MachineId)
	}

	return minResourceCostMachine
}

func (s *BestFitStrategy) bestFitCpuCost(
	instance *cloud.Instance, skip *cloud.Machine) *cloud.Machine {

	minCpuCost := math.MaxFloat64
	var minCpuCostMachine *cloud.Machine

	for _, m := range s.machineDeployList {
		if skip != nil && m.MachineId == skip.MachineId {
			continue
		}

		cost := m.GetCostWithInstance(instance)
		if cost > minCpuCost {
			continue
		}

		if m.ConstraintCheck(instance) {
			minCpuCost = cost
			minCpuCostMachine = m
		}
	}

	if minCpuCostMachine != nil {
		//fmt.Printf("BestFitStrategy.bestFitResource cost=%f,%f,machineId=%d\n",
		//	minCost, instance.Config.Cpu[0], minCostMachine.MachineId)
	}

	return minCpuCostMachine
}
