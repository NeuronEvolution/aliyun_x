package fss

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
)

func (s *FreeSmallerStrategy) calcMachineCostPlusInstance(m *cloud.Machine, instance *cloud.Instance) float64 {
	totalCost := float64(0)
	for i := 0; i < cloud.TimeSampleCount; i++ {
		s := 1 + 10*(math.Exp(math.Max(0, (m.Cpu[i]+instance.Config.Cpu[i])/m.LevelConfig.Cpu-0.5))-1)
		totalCost += s
	}

	return totalCost / cloud.TimeSampleCount
}

func (s *FreeSmallerStrategy) findFirstFit(instance *cloud.Instance, skip *cloud.Machine) *cloud.Machine {
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

func (s *FreeSmallerStrategy) findHighLevelCpuAvailable(instance *cloud.Instance, skip *cloud.Machine) (m *cloud.Machine) {
	if len(s.R.MachineDeployPool.MachineLevelDeployArray) > 0 {
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
	}

	return nil
}

func (s *FreeSmallerStrategy) findSmall(instance *cloud.Instance, skip *cloud.Machine) (m *cloud.Machine) {
	//instanceList := s.R.GetInstanceOrderByCodeDescList()
	return nil
}

func (s *FreeSmallerStrategy) findAvailableMachine(instance *cloud.Instance, skip *cloud.Machine) (m *cloud.Machine) {
	m = s.findHighLevelCpuAvailable(instance, skip)
	if m != nil {
		return m
	}

	m = s.R.MachineFreePool.PeekMachine()
	if m != nil {
		if skip != nil && m.MachineId != skip.MachineId {
			fmt.Printf("FreeSmallerStrategy.findAvailableMachine skip self machindId=%d", m.MachineId)
			return nil
		}

		if m.ConstraintCheck(instance) {
			return m
		}
	}

	//todo 资源最低fit
	panic("")
	return s.findFirstFit(instance, skip)
}
