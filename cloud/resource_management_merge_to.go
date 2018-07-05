package cloud

import (
	"fmt"
)

func (r *ResourceManagement) margeMachine(m *Machine, pool []*Machine, status *ResourceManagement) (err error) {
	instanceList := make([]*Instance, 0)
	instanceList = append(instanceList, m.InstanceArray[:m.InstanceArrayCount]...)
	for _, instance := range instanceList {
		hasFit := false
		for _, freeMachine := range pool {
			if freeMachine.ConstraintCheck(instance) {
				m.RemoveInstance(instance.InstanceId)
				r.CommandDeployInstance(instance, freeMachine)
				hasFit = true
				break
			}
		}
		if !hasFit {
			freeMachine := r.MachineFreePool.PeekMachine()
			if freeMachine == nil {
				return fmt.Errorf("ResourceManagement.MergeTo no free machine")
			}

			if !freeMachine.ConstraintCheck(instance) {
				return fmt.Errorf("ResourceManagement.MergeTo free machine ConstraintCheck failed")
			}

			m.RemoveInstance(instance.InstanceId)
			r.CommandDeployInstance(instance, freeMachine)
		}
	}

	for _, instance := range r.InstanceList {
		if instance == nil {
			continue
		}

		machineId := status.InstanceDeployedMachineMap[instance.InstanceId].MachineId
		if machineId == m.MachineId {
			if !m.ConstraintCheck(instance) {
				return fmt.Errorf("ResourceManagement.MergeTo return back ConstraintCheck failed")
			}
			oldMachine := r.InstanceDeployedMachineMap[instance.InstanceId]
			oldMachine.RemoveInstance(instance.InstanceId)
			r.CommandDeployInstance(instance, m)
		}
	}

	return nil
}

func (r *ResourceManagement) mergeNonDeployedMachines(status *ResourceManagement) (err error) {
	for _, instance := range status.InstanceList {
		if instance == nil {
			continue
		}

		if r.InstanceList[instance.InstanceId] == nil {
			newInstance := r.CreateInstance(instance.InstanceId, instance.Config)
			m := r.MachineMap[status.InstanceDeployedMachineMap[instance.InstanceId].MachineId]
			if !m.ConstraintCheck(instance) {
				return fmt.Errorf("ResourceManagement.MergeTo ConstraintCheck failed %d %d\n",
					m.MachineId, instance.InstanceId)
			}
			r.CommandDeployInstance(newInstance, m)
		}
	}

	return nil
}

func (r *ResourceManagement) MergeTo(status *ResourceManagement) (err error) {
	var mapping [MaxMachineId]int
	mappedCount := 0
	for _, instance := range r.InstanceList {
		if instance == nil {
			continue
		}

		srcMachine := r.InstanceDeployedMachineMap[instance.InstanceId]
		machineId := status.InstanceDeployedMachineMap[instance.InstanceId].MachineId
		if mapping[machineId] == 0 {
			mapping[machineId] = srcMachine.MachineId
			mappedCount++
		}
	}

	initMachineMap := make([]*Machine, 0)
	for _, level := range r.MachineDeployPool.MachineLevelDeployArray {
		initMachineMap = append(initMachineMap, level.MachineCollection.List[:level.MachineCollection.ListCount]...)
	}

	freeMachineMap := make([]*Machine, 0)
	for _, level := range r.MachineFreePool.MachineLevelFreeArray {
		freeMachineMap = append(freeMachineMap, level.MachineCollection.List[:level.MachineCollection.ListCount]...)
	}

	for _, m := range initMachineMap {
		err = r.margeMachine(m, freeMachineMap, status)
		if err != nil {
			return err
		}
	}

	fmt.Printf("ResourceManagement.MergeTo freed machine count=%d\n", len(freeMachineMap))

	for _, m := range freeMachineMap {
		err = r.margeMachine(m, initMachineMap, status)
		if err != nil {
			return err
		}
	}

	return r.mergeNonDeployedMachines(status)
}
