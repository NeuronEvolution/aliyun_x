package bfs

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
)

func (s *BestFitStrategy) bestFit(
	instance *cloud.Instance, skip *cloud.Machine, cpuMax float64) *cloud.Machine {

	minResourceCost := math.MaxFloat64
	var minResourceCostMachine *cloud.Machine

	for _, m := range s.machineDeployList {
		if skip != nil && m.MachineId == skip.MachineId {
			continue
		}

		if m.ResourceCost >= minResourceCost {
			continue
		}

		cost := m.GetCostWithInstance(instance)
		if cost > cpuMax {
			//fmt.Printf("calcMachineRealCostPlusInstance cost > cpuLimit %d %d %f\n",
			//	m.MachineId, instance.InstanceId, cost)
			continue
		}

		if m.ConstraintCheck(instance) {
			minResourceCost = m.ResourceCost
			minResourceCostMachine = m
		}
	}

	if minResourceCostMachine != nil {
		//fmt.Printf("BestFitStrategy.bestFit cost=%f,%f,machineId=%d\n",
		//	minCost, instance.Config.Cpu[0], minCostMachine.MachineId)
	}

	return minResourceCostMachine
}
