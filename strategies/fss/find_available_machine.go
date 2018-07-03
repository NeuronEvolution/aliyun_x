package fss

import (
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

func (s *FreeSmallerStrategy) bestFit(
	instance *cloud.Instance, skip *cloud.Machine, cpuMax float64) *cloud.Machine {

	minDisk := math.MaxInt64
	var minDiskMachine *cloud.Machine

	for _, m := range s.machineDeployList {
		if skip != nil && m.MachineId == skip.MachineId {
			continue
		}

		if m.Disk >= minDisk {
			continue
		}

		cost := m.GetCostWithInstance(instance)
		if cost > cpuMax {
			//fmt.Printf("calcMachineRealCostPlusInstance cost > cpuLimit %d %d %f\n",
			//	m.MachineId, instance.InstanceId, cost)
			continue
		}

		if m.ConstraintCheck(instance) {
			minDisk = m.Disk
			minDiskMachine = m
		}
	}

	if minDiskMachine != nil {
		//fmt.Printf("FreeSmallerStrategy.bestFit cost=%f,%f,machineId=%d\n",
		//	minCost, instance.Config.Cpu[0], minCostMachine.MachineId)
	}

	return minDiskMachine
}
