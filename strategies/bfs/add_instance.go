package bfs

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
	"sort"
)

func (s *BestFitStrategy) AddInstanceList(instanceList []*cloud.Instance) (err error) {
	sort.Sort(cloud.InstanceListSortByCostEvalDesc(instanceList))
	for i, v := range instanceList {
		if i > 0 && i%1000 == 0 {
			fmt.Println(i)
		}

		if i > 65000 {
			cloud.SetDebug(true)
		}

		err = s.addInstance(v, nil)
		if err != nil {
			return err
		}
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

func (s *BestFitStrategy) addInstance(instance *cloud.Instance, skip *cloud.Machine) (err error) {
	//0.6CPU内，插入后最小原则插入
	m := s.bestFit(instance, skip, cloud.MaxCpu)
	if m != nil {
		s.R.CommandDeployInstance(instance, m)
		s.sortMachineDeployList()
		return nil
	}

	//return fmt.Errorf("BestFitStrategy.addInstance no machine")

	//fmt.Printf("BestFitStrategy.addInstance no machine\n")

	m = s.bestFit(instance, skip, math.MaxFloat64)
	if m == nil {
		return fmt.Errorf("BestFitStrategy.addInstance bestFit failed")
	}

	s.R.CommandDeployInstance(instance, m)
	s.sortMachineDeployList()
	return nil
}
