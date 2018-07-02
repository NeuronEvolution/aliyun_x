package fss

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

func (s *FreeSmallerStrategy) redeployInstanceList(instanceList []*cloud.Instance, breakOnFail bool) error {
	sort.Sort(cloud.InstanceListSortByCostEvalDesc(instanceList))
	for _, v := range instanceList {
		m := s.R.InstanceDeployedMachineMap[v.InstanceId]
		m.RemoveInstance(v.InstanceId)

		newMachine := s.findAvailableMachine(v, m)
		if newMachine == nil {
			m.AddInstance(v)

			if breakOnFail {
				return fmt.Errorf("FreeSmallerStrategy.resolveAppInference findAvailableMachine failed,"+
					"machineId=%d,instanceId=%d\n",
					m.MachineId, v.InstanceId)
			} else {
				continue
			}
		}

		s.R.CommandDeployInstance(v, newMachine)
	}

	return nil
}
