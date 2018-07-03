package bfs

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

type BestFitStrategy struct {
	R                 *cloud.ResourceManagement
	machineDeployList []*cloud.Machine
}

func NewFreeSmallerStrategy(r *cloud.ResourceManagement) *BestFitStrategy {
	s := &BestFitStrategy{}
	s.R = r

	return s
}

func (s *BestFitStrategy) Name() string {
	return "BestFit"
}

func (s *BestFitStrategy) sortMachineDeployList() {
	sort.Slice(s.machineDeployList, func(i, j int) bool {
		return s.machineDeployList[i].ResourceCost < s.machineDeployList[j].ResourceCost
	})
}
