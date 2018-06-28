package clound

import "fmt"

//添加机器
//todo 将触发调度
//todo 异步将最低资源机器释放，相关instance重新拉入
func (r *ResourceManagement) AddMachine(config *MachineResourcesConfig) error {
	r.machineConfigMap[config.MachineId] = config
	machine := NewMachine(r, config.MachineId, r.machineLevelConfigPool.GetConfig(&config.MachineLevelConfig))
	r.machineMap[machine.MachineId] = machine
	r.addFreeMachine(machine)

	return nil
}

//删除机器
//不能删除有实例部署的机器
//todo 异步将最低资源机器释放，相关instance重新拉入
func (r *ResourceManagement) RemoveMachine(machineId string) error {
	m := r.machineMap[machineId]
	if m.InstanceArrayCount > 0 {
		return fmt.Errorf("ResourceManagement.RemoveMachine 机器%s已部署%d个实例",
			machineId, m.InstanceArrayCount)
	}

	delete(r.machineConfigMap, machineId)
	delete(r.machineMap, machineId)
	r.removeFromFreeMachine(machineId)

	return nil
}

func (r *ResourceManagement) addFreeMachine(m *Machine) {
	r.machineFreePool.AddMachine(m)
}

func (r *ResourceManagement) removeFromFreeMachine(machineId string) {
	r.machineFreePool.RemoveMachine(machineId)
}

func (r *ResourceManagement) popFreeMachine() (machine *Machine) {
	return r.machineFreePool.PopMachine()
}
