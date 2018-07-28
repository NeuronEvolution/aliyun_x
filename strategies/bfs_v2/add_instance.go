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

	restInstances, err := s.mem8(instances)
	if err != nil {
		return err
	}

	fmt.Println("mem8 rest count", len(restInstances))

	cloud.SortInstanceByTotalMaxLow(restInstances)

	for i, m := range s.machineDeployList {
		if i >= 3000 {
			break
		}

		fmt.Println("predploy", i, len(restInstances))
		restInstances, err = s.preDeploy(m, restInstances)
		//m.DebugPrint()
		//fmt.Println(m.Resource.GetCpuCost(m.LevelConfig.Cpu), m.Resource.GetLinearCpuCost(m.LevelConfig.Cpu))
		if err != nil {
			return err
		}
	}

	cloud.SortInstanceByTotalMaxLowWithInference(restInstances)

	fmt.Println("AddInstanceList restInstances ", len(restInstances))
	for i, v := range restInstances {
		if i < 0 {
			v.Config.DebugPrint()
		}
	}

	for i, v := range restInstances {
		//if i > 0 && i%1000 == 0 {
		fmt.Println(i)
		//}

		err = s.addInstance(v, float64(i)/float64(len(restInstances)))
		if err != nil {
			//for _, m := range s.machineDeployList {
			//m.DebugPrint()
			//}
			fmt.Println(i)
			return err
		}
	}

	sort.Slice(s.machineDeployList, func(i, j int) bool {
		if s.machineDeployList[i].LevelConfig.Cpu > s.machineDeployList[j].LevelConfig.Cpu {
			return true
		} else if s.machineDeployList[i].LevelConfig.Cpu == s.machineDeployList[j].LevelConfig.Cpu {
			return s.machineDeployList[i].GetCpuCost() > s.machineDeployList[j].GetCpuCost()
		} else {
			return false
		}
	})
	for _, m := range s.machineDeployList {
		fmt.Printf("cost=%f %f\n", m.GetCpuCost(), m.GetLinearCpuCost(m.LevelConfig.Cpu))
		if m.GetLinearCpuCost(m.LevelConfig.Cpu) < 0.9 {
			m.DebugPrint()
		} else {
			m.Resource.DebugPrint()
		}
	}

	return nil
}

func (s *Strategy) addInstance(instance *cloud.Instance, progress float64) (err error) {
	m := s.bestFitResource(instance, cloud.MaxCpuRatio, progress)
	if m != nil {
		m.AddInstance(instance)
		return nil
	}

	m = s.bestFitCpuCost(instance, progress, true)
	if m == nil {
		return fmt.Errorf("BestFitStrategy.addInstance bestFitCpuCost failed")
	}

	m.AddInstance(instance)

	return nil
}
