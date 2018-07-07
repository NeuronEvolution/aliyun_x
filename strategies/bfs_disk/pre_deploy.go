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

func (s *BestFitStrategy) preDeployDisk1000(instanceList []*cloud.Instance) (result []*cloud.Instance, err error) {
	deployed := make(map[int]*cloud.Instance)
	for _, instance := range instanceList {
		if instance.Config.Disk < 1000 {
			continue
		}

		for _, m := range s.machineDeployList {
			if m.InstanceArrayCount == 0 {
				if !m.ConstraintCheck(instance) {
					return nil, fmt.Errorf("BestFitStrategy.preDeployDisk1000 ConstraintCheck failed")
				}

				fmt.Printf("    preDeployDisk1000 instanceId=%d\n", instance.InstanceId)

				s.R.CommandDeployInstance(instance, m)
				deployed[instance.InstanceId] = instance

				break
			}
		}
	}

	fmt.Printf("BestFitStrategy.preDeployDisk1000 totalCount=%d\n", len(deployed))

	return s.preDeployRemoveDeployed(instanceList, deployed), nil
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

func (s *BestFitStrategy) preDeploy(instanceList []*cloud.Instance) (result []*cloud.Instance, err error) {
	result, err = s.preDeployDisk1000(instanceList)
	if err != nil {
		return nil, err
	}

	result, err = s.preDeployDisk167(instanceList)
	if err != nil {
		return nil, err
	}

	totalMachineCount := 0
	for _, m := range s.machineDeployList {
		if m.InstanceArrayCount > 0 {
			totalMachineCount++
		}
	}

	fmt.Printf("BestFitStrategy.preDeploy totalInstanceCount=%d,totalMachineCount=%d\n",
		len(instanceList)-len(result), totalMachineCount)

	return result, nil
}
