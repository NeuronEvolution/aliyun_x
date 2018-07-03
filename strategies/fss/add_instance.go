package fss

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
	"sort"
)

func (s *FreeSmallerStrategy) AddInstanceList(instanceList []*cloud.Instance) (err error) {
	sort.Slice(s.machineDeployList, func(i, j int) bool {
		return s.machineDeployList[i].Disk < s.machineDeployList[j].Disk
	})
	cloud.SetDebug(true)
	sort.Sort(cloud.InstanceListSortByCostEvalDesc(instanceList))
	for i, v := range instanceList {
		if i > 0 && i%1000 == 0 {
			fmt.Println(i)
		}

		err = s.addInstance(v, nil)
		if err != nil {
			return err
		}
	}

	err = s.resolveHighCpuMachine()
	if err != nil {
		return err
	}

	return nil
}

func (s *FreeSmallerStrategy) getDeployMachineList(totalCount int) []*cloud.Machine {
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

func (s *FreeSmallerStrategy) addInstance(instance *cloud.Instance, skip *cloud.Machine) (err error) {
	//0.6CPU内，插入后最小原则插入
	m := s.bestFit(instance, skip, cloud.MaxCpu)
	if m != nil {
		s.R.CommandDeployInstance(instance, m)
		sort.Slice(s.machineDeployList, func(i, j int) bool {
			return s.machineDeployList[i].Disk < s.machineDeployList[j].Disk
		})
		return nil
	}

	return fmt.Errorf("FreeSmallerStrategy.addInstance no machine")

	//fmt.Printf("FreeSmallerStrategy.addInstance no machine\n")

	m = s.bestFit(instance, skip, math.MaxFloat64)
	if m == nil {
		return fmt.Errorf("FreeSmallerStrategy.addInstance bestFit failed")
	}

	s.R.CommandDeployInstance(instance, m)
	sort.Slice(s.machineDeployList, func(i, j int) bool {
		return s.machineDeployList[i].Disk < s.machineDeployList[j].Disk
	})
	return nil
}
