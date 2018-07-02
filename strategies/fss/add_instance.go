package fss

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

func (s *FreeSmallerStrategy) AddInstance(instance *cloud.Instance) (err error) {
	m := s.findAvailableMachine(instance, nil)
	if m == nil {
		return fmt.Errorf("SortedFirstFitStrategy.AddInstance no findFirstFit")
	}

	s.R.CommandDeployInstance(instance, m)

	return nil
}

func (s *FreeSmallerStrategy) AddInstanceList(instanceList []*cloud.Instance) (err error) {
	sort.Sort(cloud.InstanceListSortByCostEvalDesc(instanceList))
	for i, v := range instanceList {
		if i > 0 && i%1000 == 0 {
			fmt.Println(i)
		}

		m := s.findAvailableMachine(v, nil)
		if m == nil {
			return fmt.Errorf("SortedFirstFitStrategy.AddInstance no findFirstFit")
		}

		s.R.CommandDeployInstance(v, m)
	}

	return nil
}
