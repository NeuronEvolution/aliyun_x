package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

func (s *BestFitStrategy) getDeployMachineList(totalCount int) []*cloud.Machine {
	machineList := make([]*cloud.Machine, 0)

	for _, v := range s.R.MachineDeployPool.MachineLevelDeployArray {
		machineList = append(machineList, v.MachineCollection.List[:v.MachineCollection.ListCount]...)
	}

	freeCount := totalCount - len(s.R.MachineDeployPool.MachineMap)
	if freeCount < 0 {
		panic(fmt.Errorf("freeCount< 0,totalCount=%d,deployed=%d\n",
			totalCount, len(s.R.MachineDeployPool.MachineMap)))
	}

	machineList = append(machineList, s.R.MachineFreePool.PeekMachineList(freeCount)...)

	return machineList
}
