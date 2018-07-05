package bfs

import (
	"fmt"
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

func (s *BestFitStrategy) PostInit() (err error) {
	s.machineDeployList = s.getDeployMachineList(MachineDeployCount)
	if len(s.machineDeployList) != MachineDeployCount {
		panic("BestFitStrategy.AddInstanceList getDeployMachineList failed")
	}

	return nil
}

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

func (s *BestFitStrategy) sortMachineDeployList() {
	sort.Slice(s.machineDeployList, func(i, j int) bool {
		return s.machineDeployList[i].ResourceCost < s.machineDeployList[j].ResourceCost
	})
}
