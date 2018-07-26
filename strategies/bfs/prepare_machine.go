package bfs

import (
	"sort"
)

func (s *BestFitStrategy) sortMachineDeployList() {
	sort.Slice(s.machineDeployList, func(i, j int) bool {
		return s.machineDeployList[i].ResourceCost < s.machineDeployList[j].ResourceCost
	})
}
