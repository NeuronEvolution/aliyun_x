package bfs_disk

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

func (s *BestFitStrategy) AddInstanceList(instanceList []*cloud.Instance) (err error) {
	s.machineDeployList = s.getDeployMachineList(MachineDeployCount)
	if len(s.machineDeployList) != MachineDeployCount {
		panic("BestFitStrategy.AddInstanceList getDeployMachineList failed")
	}

	restInstances, restMachines, err := s.preDeploy(instanceList)
	if err != nil {
		return err
	}

	s.machineDeployList = restMachines
	fmt.Printf("AddInstanceList machineCount=%d\n", len(s.machineDeployList))

	sort.Sort(cloud.InstanceListSortByCostEvalDesc(restInstances))
	for i, v := range restInstances {
		if i > 0 && i%1000 == 0 {
			fmt.Println(i)
		}

		err = s.addInstance(v, nil)
		if err != nil {
			fmt.Println(i)
			return err
		}
	}

	return nil
}

func (s *BestFitStrategy) addInstance(instance *cloud.Instance, skip *cloud.Machine) (err error) {
	//0.6CPU内，插入后最小原则插入
	m := s.bestFitResource(instance, skip, cloud.MaxCpu)
	if m != nil {
		m.AddInstance(instance)
		s.sortMachineDeployList()
		return nil
	}

	m = s.bestFitCpuCost(instance, skip)
	if m == nil {
		return fmt.Errorf("BestFitStrategy.addInstance bestFitResource failed")
	}

	m.AddInstance(instance)
	s.sortMachineDeployList()
	return nil
}
