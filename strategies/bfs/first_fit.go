package bfs

import "github.com/NeuronEvolution/aliyun_x/cloud"

func (s *BestFitStrategy) firstFit(instance *cloud.Instance, skip *cloud.Machine) *cloud.Machine {
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
