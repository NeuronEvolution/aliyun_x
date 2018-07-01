package cloud

import "fmt"

func (r *ResourceManagement) CommandAddMachine(config *MachineResourcesConfig) error {
	return r.AddMachine(config)
}

func (r *ResourceManagement) CommandRemoveMachine(machineId int) error {
	return r.RemoveMachine(machineId)
}

func (r *ResourceManagement) CommandDeployInstance(instance *Instance, m *Machine) {
	//debugLog("ResourceManagement.CommandDeployInstance appId=%d,instanceId=%d,machineId=%d",
	//	instance.Config.AppId, instance.InstanceId, m.MachineId)

	m.AddInstance(instance)

	r.DeployCommandHistory.Push(instance.Config.AppId, instance.InstanceId, m.MachineId)
}

func (r *ResourceManagement) Play(h *DeployCommandHistory) (err error) {
	debugLog("ResourceManagement.Play command count=%d", h.ListCount)
	if h == nil {
		return fmt.Errorf("ResourceManagement.Play arg nil")
	}

	for i, v := range h.List[:h.ListCount] {
		m := r.MachineMap[v.MachineId]
		if m == nil {
			return fmt.Errorf("ResourceManagement.Play machine %d not exists", v.MachineId)
		}

		appResourcesConfig := r.AppResourcesConfigMap[v.AppId]
		if appResourcesConfig == nil {
			return fmt.Errorf("ResourceManagement.Play instanceId=%d,appId=%d appResourcesConfig not found",
				v.InstanceId, v.AppId)
		}

		instance := r.InstanceList[v.InstanceId]
		if instance == nil {
			instance = NewInstance(r, v.InstanceId, appResourcesConfig)
		}

		currentMachine := r.InstanceMachineMap[instance.InstanceId]
		if currentMachine != nil {
			if currentMachine.MachineId == m.MachineId {
				panic(fmt.Errorf("ResourceManagement.Play self deploy %d %v", i, v))
			}

			currentMachine.RemoveInstance(v.InstanceId)
		}

		if !m.ConstraintCheck(instance) {
			return fmt.Errorf("ResourceManagement.Play ConstraintCheck failed %d %v ", i, v)
		}
		m.AddInstance(instance)
	}

	return nil
}