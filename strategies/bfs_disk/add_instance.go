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

	sort.Sort(cloud.InstanceListSortByCostEvalDesc(instanceList))
	for i, v := range instanceList {
		if i > 0 && i%1000 == 0 {
			fmt.Println(i)
		}

		if i > 68000 {
			cloud.SetDebug(true)
		}

		err = s.addInstance(v, nil)
		if err != nil {
			fmt.Println(i)
			return err
		}
	}

	cloud.SetDebug(false)

	return nil
}

func (s *BestFitStrategy) addInstance(instance *cloud.Instance, skip *cloud.Machine) (err error) {
	//0.6CPU内，插入后最小原则插入
	m := s.bestFitResource(instance, skip, cloud.MaxCpu)
	if m != nil {
		s.R.CommandDeployInstance(instance, m)
		s.sortMachineDeployList()
		return nil
	}

	m = s.bestFitCpuCost(instance, skip)
	if m == nil {
		return fmt.Errorf("BestFitStrategy.addInstance bestFitResource failed")
	}

	s.R.CommandDeployInstance(instance, m)
	s.sortMachineDeployList()
	return nil
}
