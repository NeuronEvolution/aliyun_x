package cloud

import (
	"fmt"
	"sort"
)

//该操作将重置所有实例部署
//todo 异步化
func (r *ResourceManagement) InitInstanceDeploy(configList []*InstanceDeployConfig) error {
	if configList == nil || len(configList) == 0 {
		return nil
	}

	r.InitialInstanceDeployConfig = configList

	for _, v := range configList {
		appResourcesConfig := r.AppResourcesConfigMap[v.AppId]
		if appResourcesConfig == nil {
			return fmt.Errorf("R.InitInstanceDeploy %d appResourcesConfig %d not found",
				v.InstanceId, v.AppId)
		}
		instance := NewInstance(r, v.InstanceId, appResourcesConfig)

		m := r.MachineMap[v.MachineId]
		if m == nil {
			return fmt.Errorf("ResourceManagement.InitInstanceDeploy machine %d not exsits", v.MachineId)
		}

		if !constraintCheckResourceLimit(m, instance) {
			//debugLog("Machine.ConstraintCheck constraintCheckResourceLimit failed")
			return fmt.Errorf("constraintCheckResourceLimit failed %d", m.MachineId)
		}

		m.AddInstance(instance)
	}

	return nil
}

//todo 异步化
func (r *ResourceManagement) AddInstance(c *InstanceDeployConfig) error {
	debugLog("R.AddInstance %s %s", c.InstanceId, c.AppId)

	appResourcesConfig := r.AppResourcesConfigMap[c.AppId]
	if appResourcesConfig == nil {
		return fmt.Errorf("R.AddInstance %d appResourcesConfig %d not found",
			c.InstanceId, c.AppId)
	}

	instance := NewInstance(r, c.InstanceId, appResourcesConfig)

	return r.Strategy.AddInstance(instance)
}

//todo 异步化
func (r *ResourceManagement) AddInstanceList(configList []*InstanceDeployConfig) error {
	if configList == nil || len(configList) == 0 {
		return nil
	}

	instanceList := make(InstanceListSortByCostEvalDesc, 0)
	for _, c := range configList {
		appResourcesConfig := r.AppResourcesConfigMap[c.AppId]
		if appResourcesConfig == nil {
			return fmt.Errorf("R.AddInstance %d appResourcesConfig %d not found",
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

func (r *ResourceManagement) SetInstanceDeployedMachine(instance *Instance, m *Machine) {
	if r.InstanceDeployedMachineMap[instance.InstanceId] == nil {
		r.instanceDeployedOrderByCostDescValid = false
	}

	r.InstanceDeployedMachineMap[instance.InstanceId] = m
}

func (r *ResourceManagement) GetInstanceOrderByCodeDescList() (instanceList []*Instance) {
	if !r.instanceDeployedOrderByCostDescValid {
		r.instanceDeployedOrderByCostDescListCount = 0
		for i := range r.InstanceDeployedMachineMap {
			r.instanceDeployedOrderByCostDescList[r.instanceDeployedOrderByCostDescListCount] = r.InstanceList[i]
			r.instanceDeployedOrderByCostDescListCount++
		}
		sort.Sort(InstanceListSortByCostEvalDesc(
			r.instanceDeployedOrderByCostDescList[:r.instanceDeployedOrderByCostDescListCount]))

		r.instanceDeployedOrderByCostDescValid = true
	}

	return r.instanceDeployedOrderByCostDescList[:r.instanceDeployedOrderByCostDescListCount]
}
