package strategies

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

type AllocMachineIfDeployFailedStrategy struct {
	R *cloud.ResourceManagement
}

func NewAllocMachineIfDeployFailedStrategy(r *cloud.ResourceManagement) cloud.Strategy {
	s := &AllocMachineIfDeployFailedStrategy{}
	s.R = r

	return s
}

func (s *AllocMachineIfDeployFailedStrategy) AddInstance(instance *cloud.Instance) (err error) {
	for _, v := range s.R.MachineDeployPool.MachineLevelDeployArray {
		for i := 0; i < v.MachineCollection.ListCount; i++ {
			m := v.MachineCollection.List[i]
			if m.ConstraintCheck(instance) {
				m.AddInstance(instance)
				return nil
			}
		}
	}

	m := s.R.MachineFreePool.PopMachine()
	if m == nil {
		return fmt.Errorf("no free machine")
	}

	if !m.ConstraintCheck(instance) {
		return fmt.Errorf("AllocMachineIfDeployFailedStrategy.AddInstance ConstraintCheck failed")
	}
	m.AddInstance(instance)
	s.R.MachineDeployPool.AddMachine(m)

	return nil
}

func (s *AllocMachineIfDeployFailedStrategy) AddInstanceList(instanceList []*cloud.Instance) (err error) {
	sort.Sort(cloud.InstanceArray(instanceList))

	for i, v := range instanceList {
		//fmt.Println(v.CostEval)

		if i%1000 == 0 {
			fmt.Println(i)
		}

		err = s.AddInstance(v)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return
}
