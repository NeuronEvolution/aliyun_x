package clound

import "fmt"

//该操作将重置所有实例部署
//todo 异步化
func (r *ResourceManagement) InitInstanceDeploy(list []*InstanceDeployConfig) error {
	return nil
}

//todo 异步化
func (r *ResourceManagement) AddInstance(instanceId string, appId string) error {
	debugLog("ResourceManagement.AddInstance %s %s", instanceId, appId)

	appResourcesConfig := r.appResourcesConfigMap[appId]
	if appResourcesConfig == nil {
		return fmt.Errorf("ResourceManagement.AddInstance %s appResourcesConfig %s not found",
			instanceId, appId)
	}

	instance := NewInstance(r, instanceId, appResourcesConfig)

	var machine *Machine
	for _, v := range r.machineDeployPool.MachineMap {
		machine = v
		break
	}
	if machine == nil {
		machine = r.machineFreePool.PopMachine()
	}

	if machine == nil {
		return fmt.Errorf("no free machine")
	}
	err := machine.AddInstance(instance)
	if err != nil {
		return err
	}

	r.machineDeployPool.AddMachine(machine)

	return nil
}

//todo 异步化
func (r *ResourceManagement) RemoveInstance(instanceId string) error {
	return nil
}
