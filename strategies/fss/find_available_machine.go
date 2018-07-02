package fss

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
)

func (s *FreeSmallerStrategy) calcMachineRealCostPlusInstance(m *cloud.Machine, instance *cloud.Instance) float64 {
	totalCost := float64(0)
	for i := 0; i < cloud.TimeSampleCount; i++ {
		r := (m.Cpu[i] + instance.Config.Cpu[i]) / m.LevelConfig.Cpu
		s := 1 + 10*(math.Exp(math.Max(0, r-0.5))-1) - (0.5 - math.Min(0.5, r))
		totalCost += s
	}

	return totalCost / cloud.TimeSampleCount
}

func (s *FreeSmallerStrategy) firstFit(instance *cloud.Instance, skip *cloud.Machine) *cloud.Machine {
	for _, v := range s.R.MachineDeployPool.MachineLevelDeployArray {
		for _, m := range v.MachineCollection.List[:v.MachineCollection.ListCount] {
			if skip != nil && m.MachineId == skip.MachineId {
				continue
			}

			if m.ConstraintCheck(instance) {
				return m
			}
		}
	}

	return nil
}

func (s *FreeSmallerStrategy) bestFit(instance *cloud.Instance, skip *cloud.Machine, cpuLimit float64) *cloud.Machine {
	minCost := float64(math.MaxFloat64)
	var minCostMachine *cloud.Machine
	for _, v := range s.R.MachineDeployPool.MachineLevelDeployArray {
		for _, m := range v.MachineCollection.List[:v.MachineCollection.ListCount] {
			if skip != nil && m.MachineId == skip.MachineId {
				continue
			}

			cost := s.calcMachineRealCostPlusInstance(m, instance)
			if cost < cpuLimit && cost >= minCost {
				continue
			}

			if m.ConstraintCheck(instance) {
				minCost = cost
				minCostMachine = m
			}
		}
	}

	if minCostMachine != nil {
		//fmt.Printf("FreeSmallerStrategy.bestFit cost=%f,%f,machineId=%d\n",
		//	minCost, instance.Config.Cpu[0], minCostMachine.MachineId)
	}

	return minCostMachine
}
