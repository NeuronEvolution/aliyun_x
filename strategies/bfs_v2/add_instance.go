package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

func (s *Strategy) AddInstanceList(instances []*cloud.Instance) (err error) {
	s.machineDeployList = s.R.MachineFreePool.PeekMachineList(MachineDeployCount)
	if len(s.machineDeployList) != MachineDeployCount {
		panic("BestFitStrategy.AddInstanceList getDeployMachineList failed")
	}

	restInstances := instances
	sort.Slice(restInstances, func(i, j int) bool {
		return restInstances[i].Config.GetCpuDerivation() > restInstances[j].Config.GetCpuDerivation()
	})

	for i, m := range s.machineDeployList {
		if i >= 3000 {
			break
		}

		fmt.Println("predploy", i, len(restInstances))
		restInstances, err = s.preDeploy(m, restInstances)
		m.DebugPrint()
		fmt.Println(m.Resource.GetCpuCost(m.LevelConfig.Cpu), m.Resource.GetLinearCpuCost(m.LevelConfig.Cpu))
		if err != nil {
			return err
		}
	}

	cloud.SortInstanceByTotalMax(restInstances)

	fmt.Println("AddInstanceList restInstances ", len(restInstances))

	for i, v := range restInstances {
		//if i > 0 && i%1000 == 0 {
		fmt.Println(i)
		//}

		err = s.addInstance(v, float64(i)/float64(len(restInstances)))
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

func (s *Strategy) addInstance(instance *cloud.Instance, progress float64) (err error) {
	m := s.bestFitResource(instance, cloud.MaxCpuRatio, progress)
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
