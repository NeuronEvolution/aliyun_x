package bfs_disk

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

func (s *BestFitStrategy) preDeployRemoveDeployed(instances []*cloud.Instance, deployed map[int]*cloud.Instance) (result []*cloud.Instance) {
	result = make([]*cloud.Instance, 0)
	for _, v := range instances {
		_, has := deployed[v.InstanceId]
		if has {
			continue
		}

		result = append(result, v)
	}

	return result
}

func (s *BestFitStrategy) preDeploy(instanceList []*cloud.Instance) (
	restInstances []*cloud.Instance, restMachines []*cloud.Machine, err error) {
	restInstances, restMachines, err = s.preDeployBigDisk(instanceList, s.machineDeployList)
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("BestFitStrategy.preDeploy totalInstanceCount=%d,totalMachineCount=%d\n",
		len(instanceList)-len(restInstances), len(s.machineDeployList)-len(restMachines))

	return restInstances, restMachines, nil
}
