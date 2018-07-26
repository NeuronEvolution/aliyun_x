package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

func (s *Strategy) AddInstanceList(instanceList []*cloud.Instance) (err error) {
	s.machineDeployList = s.R.MachineFreePool.PeekMachineList(MachineDeployCount)
	if len(s.machineDeployList) != MachineDeployCount {
		panic("BestFitStrategy.AddInstanceList getDeployMachineList failed")
	}

	cloud.SortInstanceByTotalMax(instanceList)

	for i, v := range instanceList {
		//if i > 0 && i%1000 == 0 {
		fmt.Println(i)
		//}

		err = s.addInstance(v)
		if err != nil {
			for _, m := range s.machineDeployList {
				m.Resource.DebugPrint()
			}
			fmt.Println(i)
			return err
		}
	}

	for _, m := range s.machineDeployList {
		m.Resource.DebugPrint()
	}

	return nil
}

func (s *Strategy) addInstance(instance *cloud.Instance) (err error) {
	m := s.bestFitResource(instance, cloud.MaxCpuRatio)
	if m != nil {
		m.AddInstance(instance)
		return nil
	}

	m = s.bestFitCpuCost(instance)
	if m == nil {
		return fmt.Errorf("BestFitStrategy.addInstance bestFitCpuCost failed")
	}

	m.AddInstance(instance)
	//fmt.Printf("cpu ")
	//m.Resource.DebugPrint()

	return nil
}
