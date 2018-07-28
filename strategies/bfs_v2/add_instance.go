package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
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

		if i > 0 && i%100 == 0 {
			fmt.Println("predploy", i, len(restInstances))
		}
		restInstances, err = s.preDeploy(m, restInstances)
		//fmt.Println(m.Resource.GetCpuCost(m.LevelConfig.Cpu), m.Resource.GetLinearCpuCost(m.LevelConfig.Cpu))
		if err != nil {
			return err
		}
	}

	cloud.SortInstanceByTotalMaxLowWithInference(restInstances)

	fmt.Println("AddInstanceList restInstances ", len(restInstances))
	for i, v := range restInstances {
		if i > 0 && i%1000 == 0 {
			fmt.Println(i)
		}

		err = s.addInstance(v, float64(i)/float64(len(restInstances)))
		if err != nil {
			fmt.Println(i)
			return err
		}
	}

	s.debug()

	return nil
}

func (s *Strategy) addInstance(instance *cloud.Instance, progress float64) (err error) {
	m := s.bestFitCpuCost(instance, progress, false)
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
