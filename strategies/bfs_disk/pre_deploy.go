package bfs_disk

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

func (s *BestFitStrategy) preDeployRemoveDeployed(
	instanceList []*cloud.Instance, deployed map[int]*cloud.Instance) (result []*cloud.Instance) {

	result = make([]*cloud.Instance, 0)
	for _, v := range instanceList {
		_, has := deployed[v.InstanceId]
		if has {
			continue
		}

		result = append(result, v)
	}

	return result
}

func (s *BestFitStrategy) preDeployDisk167(instanceList []*cloud.Instance) (result []*cloud.Instance, err error) {
	const diskSize = 167
	deployed := make(map[int]*cloud.Instance)
	for _, instance := range instanceList {
		if instance.Config.Disk != diskSize {
			continue
		}

		for _, m := range s.machineDeployList {
			if m.InstanceArrayCount == 0 || (m.InstanceArrayCount == 1 && m.InstanceArray[0].Config.Disk == diskSize) {
				if !m.ConstraintCheck(instance) {
					return nil, fmt.Errorf("BestFitStrategy.preDeployDisk167 ConstraintCheck failed")
				}

				fmt.Printf("    preDeployDisk167 instanceId=%d\n", instance.InstanceId)

				s.R.CommandDeployInstance(instance, m)
				deployed[instance.InstanceId] = instance

				break
			}
		}
	}

	fmt.Printf("BestFitStrategy.preDeployDisk167 totalCount=%d\n", len(deployed))

	return s.preDeployRemoveDeployed(instanceList, deployed), nil
}

func (s *BestFitStrategy) preDeploy(instanceList []*cloud.Instance) (
	restInstances []*cloud.Instance, restMachines []*cloud.Machine, err error) {
	restInstances, restMachines, err = s.preDeployBigDisk(instanceList)
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("BestFitStrategy.preDeploy totalInstanceCount=%d,totalMachineCount=%d\n",
		len(instanceList)-len(restInstances), len(s.machineDeployList)-len(restMachines))

	return restInstances, restMachines, nil
}
