package ffs

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

type FirstFitStrategy struct {
	R *cloud.ResourceManagement
}

func NewStrategy(r *cloud.ResourceManagement) cloud.Strategy {
	s := &FirstFitStrategy{}
	s.R = r

	return s
}

func (s *FirstFitStrategy) Name() string {
	return "SortedFirstFitStrategy"
}

func (s *FirstFitStrategy) PostInit() (err error) {
	//fmt.Printf("FirstFitStrategy.PostInit\n")
	for i := 0; ; i++ {
		//fmt.Printf("FirstFitStrategy.PostInit %d\n", i)
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

func (s *FirstFitStrategy) firstFit(instance *cloud.Instance) *cloud.Machine {
	for _, v := range s.R.MachineDeployPool.MachineLevelDeployArray {
		for i := 0; i < v.MachineCollection.ListCount; i++ {
			m := v.MachineCollection.List[i]
			if m.ConstraintCheck(instance, 1) {
				return m
			}
		}
	}

	return nil
}

func (s *FirstFitStrategy) findAvailableMachine(instance *cloud.Instance) *cloud.Machine {
	m := s.firstFit(instance)
	if m != nil {
		return m
	}

	m = s.R.MachineFreePool.PeekMachine()
	if m == nil {
		fmt.Printf("FirstFitStrategy.firstFit no machine\n")
		return nil
	}

	if !m.ConstraintCheck(instance, 1) {
		fmt.Printf("FirstFitStrategy.firstFit ConstraintCheck failed machindId=%d,instanceId=%d\n",
			m.MachineId, instance.InstanceId)
		return nil
	}

	return m
}

func (s *FirstFitStrategy) AddInstance(instance *cloud.Instance) (err error) {

	return nil
}

func (s *FirstFitStrategy) AddInstanceList(instanceList []*cloud.Instance) (err error) {
	cloud.SortInstanceByTotalMax(instanceList)

	for i, v := range instanceList {
		//fmt.Println(v.CostEval)

		if i > 0 && i%1000 == 0 {
			fmt.Println(i)
		}

		m := s.findAvailableMachine(v)
		if m == nil {
			return fmt.Errorf("FirstFitStrategy.AddInstance no firstFit")
		}

		s.R.CommandDeployInstance(v, m)
	}

	return
}
