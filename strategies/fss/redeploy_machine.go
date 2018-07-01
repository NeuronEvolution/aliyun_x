package fss

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

func (s *FreeSmallerStrategy) redeployMachine(m *cloud.Machine, breakOnFail bool) error {
	instanceList := make([]*cloud.Instance, m.InstanceArrayCount)
	for index, v := range m.InstanceArray[:m.InstanceArrayCount] {
		instanceList[index] = v
	}
	for _, v := range instanceList {
		m.RemoveInstance(v.InstanceId)

		newMachine := s.findAvailableMachine(v, m)
		if newMachine == nil {
			m.AddInstance(v)

			if breakOnFail {
				return fmt.Errorf("FreeSmallerStrategy.redeployMachine findAvailableMachine none,instanceId=%d\n",
					v.InstanceId)
			} else {
				continue
			}
		}

		s.R.CommandDeployInstance(v, newMachine)
		//fmt.Printf("redeployMachine %f,%d\n", newMachine.GetCost(), newMachine.MachineId)
	}

	return nil
}
