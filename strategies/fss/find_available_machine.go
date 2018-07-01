package fss

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
)

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

func (s *FreeSmallerStrategy) calcMachineCostPlusInstance(m *cloud.Machine, instance *cloud.Instance) float64 {
	totalCost := float64(0)
	for i := 0; i < cloud.TimeSampleCount; i++ {
		s := 1 + 10*(math.Exp(math.Max(0, (m.Cpu[i]+instance.Config.Cpu[i])/m.LevelConfig.Cpu-0.5))-1)
		totalCost += s
	}

	return totalCost / cloud.TimeSampleCount
}

func (s *FreeSmallerStrategy) findAvailableMachine(instance *cloud.Instance, skip *cloud.Machine) (m *cloud.Machine) {
	highLevelDeploy := s.R.MachineDeployPool.MachineLevelDeployArray[0]
	for _, m := range highLevelDeploy.MachineCollection.List[:highLevelDeploy.MachineCollection.ListCount] {
		if skip != nil && m.MachineId == skip.MachineId {
			continue
		}

		cost := s.calcMachineCostPlusInstance(m, instance)
		if cost > HighLevelCpuMax {
			continue
		}

		if m.ConstraintCheck(instance) {
			return m
		}
	}

	highLevelFree := s.R.MachineFreePool.MachineLevelFreeArray[0]
	if highLevelFree.MachineCollection.ListCount > 0 {
		m = s.R.MachineFreePool.PeekMachine()
		if m == nil {
			panic(fmt.Errorf("FreeSmallerStrategy.findAvailableMachine PeekMachine failed"))
		}

		if skip != nil && m.MachineId == skip.MachineId {
			return nil
		}

		if !m.ConstraintCheck(instance) {
			fmt.Printf("SortedFirstFitStrategy.firstFit ConstraintCheck failed machindId=%d,instanceId=%d\n",
				m.MachineId, instance.InstanceId)
			return nil
		}

		return m
	}

	m = s.firstFit(instance, skip)
	if m != nil {
		return m
	}

	m = s.R.MachineFreePool.PeekMachine()
	if m == nil {
		fmt.Printf("FreeSmallerStrategy.findAvailableMachine no machine\n")
		return nil
	}

	if skip != nil && m.MachineId == skip.MachineId {
		return nil
	}

	if !m.ConstraintCheck(instance) {
		fmt.Printf("FreeSmallerStrategy.findAvailableMachine ConstraintCheck"+
			" failed machindId=%d,instanceId=%d\n",
			m.MachineId, instance.InstanceId)
		return nil
	}

	return m
}
