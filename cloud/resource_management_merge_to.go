package cloud

import (
	"fmt"
)

func (r *ResourceManagement) mergeGetMachineMapping(status *ResourceManagement) (mapping []int) {
	mapping = make([]int, MaxMachineId)
	mappedCount := 0
	for _, instance := range r.InstanceList {
		if instance == nil {
			continue
		}

		srcMachine := r.InstanceDeployedMachineMap[instance.InstanceId]
		destMachineId := status.InstanceDeployedMachineMap[instance.InstanceId].MachineId
		if mapping[destMachineId] == 0 {
			srcMapped := false
			for _, v := range mapping {
				if v == srcMachine.MachineId {
					srcMapped = true
					break
				}
			}
			if srcMapped {
				continue
			}

			mapping[destMachineId] = srcMachine.MachineId
			mappedCount++
		}
	}

	fmt.Printf("ResourceManagement.mergeGetMachineMapping mappedCount=%d\n", mappedCount)

	return mapping
}

func (r *ResourceManagement) mergeGetTargetMachineId(instanceId int, status *ResourceManagement, mapping []int) int {
	targetMachine := status.InstanceDeployedMachineMap[instanceId]
	targetMachineId := targetMachine.MachineId
	if mapping[targetMachineId] != 0 {
		targetMachineId = mapping[targetMachineId]
	}

	return status.InstanceDeployedMachineMap[instanceId].MachineId
}

func (r *ResourceManagement) mergeMachine(m *Machine, pool []*Machine, status *ResourceManagement, mapping []int) (err error) {
	instanceList := make([]*Instance, 0)
	instanceList = append(instanceList, m.InstanceArray[:m.InstanceArrayCount]...)
	for _, instance := range instanceList {
		targetMachineId := r.mergeGetTargetMachineId(instance.InstanceId, status, mapping)
		if targetMachineId == m.MachineId {
			//fmt.Printf("mergeMachine self skip %d %d,%d\n",
			//	instance.InstanceId, status.InstanceDeployedMachineMap[instance.InstanceId].MachineId, m.MachineId)
			continue
		}

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
				return fmt.Errorf("ResourceManagement.mergeMachine no free machine")
			}

			if !freeMachine.ConstraintCheck(instance) {
				return fmt.Errorf("ResourceManagement.mergeMachine free machine ConstraintCheck failed")
			}

			m.RemoveInstance(instance.InstanceId)
			r.CommandDeployInstance(instance, freeMachine)
		}
	}

	for _, instance := range r.InstanceList {
		if instance == nil {
			continue
		}

		if r.InstanceDeployedMachineMap[instance.InstanceId].MachineId == m.MachineId {
			//fmt.Printf("mergeMachine self back skip %d %d\n", instance.InstanceId, m.MachineId)
			continue
		}

		targetMachineId := r.mergeGetTargetMachineId(instance.InstanceId, status, mapping)
		if targetMachineId == m.MachineId {
			if !m.ConstraintCheck(instance) {
				return fmt.Errorf("ResourceManagement.mergeMachine return back ConstraintCheck failed %d %d %d",
					m.MachineId, instance.InstanceId, status.InstanceDeployedMachineMap[instance.InstanceId].MachineId)
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
	mapping := r.mergeGetMachineMapping(status)

	initMachineMap := make([]*Machine, 0)
	for _, level := range r.MachineDeployPool.MachineLevelDeployArray {
		initMachineMap = append(initMachineMap, level.MachineCollection.List[:level.MachineCollection.ListCount]...)
	}

	freeMachineMap := make([]*Machine, 0)
	for _, level := range r.MachineFreePool.MachineLevelFreeArray {
		freeMachineMap = append(freeMachineMap, level.MachineCollection.List[:level.MachineCollection.ListCount]...)
	}

	for _, m := range initMachineMap {
		err = r.mergeMachine(m, freeMachineMap, status, mapping)
		if err != nil {
			return err
		}
	}

	fmt.Printf("ResourceManagement.MergeTo freed machine count=%d\n", len(freeMachineMap))

	for _, m := range freeMachineMap {
		err = r.mergeMachine(m, initMachineMap, status, mapping)
		if err != nil {
			return err
		}
	}

	return r.mergeNonDeployedMachines(status)
}
