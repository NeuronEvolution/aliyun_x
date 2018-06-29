package cloud

import "fmt"

//该操作将重置所有实例部署
//todo 异步化
func (r *ResourceManagement) InitInstanceDeploy(configList []*InstanceDeployConfig) error {
	if configList == nil || len(configList) == 0 {
		return nil
	}

	for _, v := range configList {
		appResourcesConfig := r.AppResourcesConfigMap[v.AppId]
		if appResourcesConfig == nil {
			return fmt.Errorf("R.InitInstanceDeploy %s appResourcesConfig %s not found",
				v.InstanceId, v.AppId)
		}
		instance := NewInstance(r, v.InstanceId, appResourcesConfig)

		m := r.MachineFreePool.RemoveMachine(v.MachineId)
		if m == nil {
			m = r.MachineDeployPool.MachineMap[v.MachineId]
			if m == nil {
				return fmt.Errorf("R.InitInstanceDeploy %s not exsits", v.MachineId)
			}
		} else {
			r.MachineDeployPool.AddMachine(m)
		}

		err := m.AddInstance(instance)
		if err != nil {
			return err
		}
	}

	return nil
}

//todo 异步化
func (r *ResourceManagement) AddInstance(c *InstanceDeployConfig) error {
	debugLog("R.AddInstance %s %s", c.InstanceId, c.AppId)

	appResourcesConfig := r.AppResourcesConfigMap[c.AppId]
	if appResourcesConfig == nil {
		return fmt.Errorf("R.AddInstance %s appResourcesConfig %s not found",
			c.InstanceId, c.AppId)
	}

	instance := NewInstance(r, c.InstanceId, appResourcesConfig)

	return r.Strategy.AddInstance(instance)
}

//todo 异步化
func (r *ResourceManagement) BatchAddInstance(configList []*InstanceDeployConfig) error {
	if configList == nil || len(configList) == 0 {
		return nil
	}

	instanceList := make(InstanceArray, 0)
	for _, c := range configList {
		appResourcesConfig := r.AppResourcesConfigMap[c.AppId]
		if appResourcesConfig == nil {
			return fmt.Errorf("R.AddInstance %s appResourcesConfig %s not found",
				c.InstanceId, c.AppId)
		}
		instance := NewInstance(r, c.InstanceId, appResourcesConfig)
		instanceList = append(instanceList, instance)
	}

	return r.Strategy.AddInstanceList(instanceList)
}

//todo 异步化
func (r *ResourceManagement) RemoveInstance(instanceId string) error {
	return nil
}
