package bfs

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

func (s *BestFitStrategy) redeployMachine(m *cloud.Machine, breakOnFail bool) error {
	instanceList := make([]*cloud.Instance, m.InstanceArrayCount)
	for index, v := range m.InstanceArray[:m.InstanceArrayCount] {
		instanceList[index] = v
	}
	for _, v := range instanceList {
		m.RemoveInstance(v.InstanceId)

		err := s.addInstance(v, m)
		if err != nil {
			m.AddInstance(v)

			if breakOnFail {
				return fmt.Errorf("BestFitStrategy.redeployMachine findAvailableMachine none,instanceId=%d\n",
					v.InstanceId)
			} else {
				continue
			}
		}
	}

	return nil
}

func (s *BestFitStrategy) redeployInstanceList(instanceList []*cloud.Instance, breakOnFail bool) error {
	sort.Sort(cloud.InstanceListSortByCostEvalDesc(instanceList))
	for _, v := range instanceList {
		m := s.R.InstanceDeployedMachineMap[v.InstanceId]
		m.RemoveInstance(v.InstanceId)

		err := s.addInstance(v, m)
		if err != nil {
			m.AddInstance(v)

			if breakOnFail {
				return fmt.Errorf("BestFitStrategy.resolveAppInference findAvailableMachine failed,"+
					"machineId=%d,instanceId=%d\n",
					m.MachineId, v.InstanceId)
			} else {
				continue
			}
		}
	}

	return nil
}
