package fss

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
	"sort"
)

func (s *FreeSmallerStrategy) rebalancedWithNewMachine(instance *cloud.Instance, m *cloud.Machine) bool {
	s.R.CommandDeployInstance(instance, m)

	return true
}

func (s *FreeSmallerStrategy) AddInstanceList(instanceList []*cloud.Instance) (err error) {
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

	return nil
}

func (s *FreeSmallerStrategy) addInstance(instance *cloud.Instance, skip *cloud.Machine) (err error) {
	//0.6CPU内，插入后最小原则插入
	m := s.bestFit(instance, skip, HighLevelCpuMax)
	if m != nil {
		s.R.CommandDeployInstance(instance, m)
		return nil
	}

	//分配新机器，重新平衡
	m = s.R.MachineFreePool.PeekMachine()
	if m != nil {
		if skip != nil && skip.MachineId == m.MachineId {
			return fmt.Errorf("FreeSmallerStrategy.addInstance skipped  ")
		}

		if s.rebalancedWithNewMachine(instance, m) {
			return nil
		}
	}

	if m == nil {
		fmt.Printf("FreeSmallerStrategy.addInstance no machine")
	} else {
		fmt.Printf("FreeSmallerStrategy.addInstance rebalancedWithNewMachine failed")
	}

	//无机器或者重新平衡失败，插入后最小原则插入
	m = s.bestFit(instance, skip, math.MaxFloat64)
	if m == nil {
		return fmt.Errorf("FreeSmallerStrategy.addInstance bestFit failed")
	}
	s.R.CommandDeployInstance(instance, m)
	return nil
}
