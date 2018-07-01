package sffs

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

type SortedFirstFitStrategy struct {
	R *cloud.ResourceManagement
}

func NewASortedFirstFitStrategy(r *cloud.ResourceManagement) cloud.Strategy {
	s := &SortedFirstFitStrategy{}
	s.R = r

	return s
}

func (s *SortedFirstFitStrategy) Name() string {
	return "SortedFirstFitStrategy"
}

func (s *SortedFirstFitStrategy) PostInit() (err error) {
	//fmt.Printf("AllocMachineIfDeployFailedStrategy.PostInit\n")
	for i := 0; ; i++ {
		//fmt.Printf("SortedFirstFitStrategy.PostInit %d\n", i)
		var m *cloud.Machine
		for _, level := range s.R.MachineDeployPool.MachineLevelDeployArray {
			for _, v := range level.MachineCollection.List[:level.MachineCollection.ListCount] {
				if v.HasBadConstraint() {
					m = v
					break
				}
			}
			if m != nil {
				break
			}
		}

		if m == nil {
			fmt.Printf("SortedFirstFitStrategy.PostInit total expand %d\n", i)
			break
		}

		instanceList := make([]*cloud.Instance, m.InstanceArrayCount)
		for index, v := range m.InstanceArray[:m.InstanceArrayCount] {
			instanceList[index] = v
		}
		for _, v := range instanceList {
			m.RemoveInstance(v.InstanceId)

			newMachine := s.findAvailableMachine(v)
			if newMachine == nil {
				return fmt.Errorf("SortedFirstFitStrategy.PostInit firstFit none,instanceId=%d\n",
					v.InstanceId)
			}
			if newMachine.MachineId == m.MachineId {
				m.AddInstance(v)
				continue
			}

			s.R.CommandDeployInstance(v, newMachine)
		}
	}

	return nil
}

func (s *SortedFirstFitStrategy) firstFit(instance *cloud.Instance) *cloud.Machine {
	for _, v := range s.R.MachineDeployPool.MachineLevelDeployArray {
		for i := 0; i < v.MachineCollection.ListCount; i++ {
			m := v.MachineCollection.List[i]
			if m.ConstraintCheck(instance) {
				return m
			}
		}
	}

	return nil
}

func (s *SortedFirstFitStrategy) findAvailableMachine(instance *cloud.Instance) *cloud.Machine {
	m := s.firstFit(instance)
	if m != nil {
		return m
	}

	m = s.R.MachineFreePool.PeekMachine()
	if m == nil {
		fmt.Printf("SortedFirstFitStrategy.firstFit no machine\n")
		return nil
	}

	if !m.ConstraintCheck(instance) {
		fmt.Printf("SortedFirstFitStrategy.firstFit ConstraintCheck failed machindId=%d,instanceId=%d\n",
			m.MachineId, instance.InstanceId)
		return nil
	}

	return m
}

func (s *SortedFirstFitStrategy) AddInstance(instance *cloud.Instance) (err error) {
	m := s.findAvailableMachine(instance)
	if m == nil {
		return fmt.Errorf("SortedFirstFitStrategy.AddInstance no firstFit")
	}

	s.R.CommandDeployInstance(instance, m)

	return nil
}

func (s *SortedFirstFitStrategy) AddInstanceList(instanceList []*cloud.Instance) (err error) {
	sort.Sort(cloud.InstanceListSortByCostEvalDesc(instanceList))

	for i, v := range instanceList {
		//fmt.Println(v.CostEval)

		if i > 0 && i%1000 == 0 {
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
