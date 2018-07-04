package cloud

import "fmt"

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

	fmt.Printf("ResourceManagement.MergeTo mappedCount=%d\n", mappedCount)
	locked := make(map[int]*Machine)
	for i := 0; ; i++ {
		fmt.Printf("ResourceManagement.MergeTo loop %d\n", i)
		hasSkipped := false
		hasMoved := false
		for _, instance := range r.InstanceList {
			if instance == nil {
				continue
			}

			srcMachine := r.InstanceDeployedMachineMap[instance.InstanceId]
			_, hasLocked := locked[srcMachine.MachineId]
			if hasLocked {
				continue
			}

			machineId := status.InstanceDeployedMachineMap[instance.InstanceId].MachineId
			mappedMachineId := machineId
			if mapping[machineId] != 0 {
				mappedMachineId = mapping[machineId]
			}
			if srcMachine.MachineId != mappedMachineId {
				targetMachine := r.MachineMap[mappedMachineId]
				if targetMachine.ConstraintCheck(instance) {
					r.CommandDeployInstance(instance, targetMachine)
					hasMoved = true
				} else {
					hasSkipped = true
				}
			}
		}

		if !hasMoved {
			var m *Machine
			for _, level := range r.MachineFreePool.MachineLevelFreeArray {
				if level.MachineCollection.ListCount > 0 {
					for _, v := range level.MachineCollection.List[:level.MachineCollection.ListCount] {
						used := false
						for _, mapped := range mapping {
							if mapped == v.MachineId {
								used = true
								break
							}
						}
						if !used {
							m = v
							break
						}
					}
				}

				if m != nil {
					break
				}
			}

			if m == nil {
				return fmt.Errorf("ResourceManagement.MergeTo no free machine")
			}

			m = r.MachineFreePool.RemoveMachine(m.MachineId)
			locked[m.MachineId] = m
			fmt.Printf("ResourceManagement.MergeTo lockedCount=%d", len(locked))
			if len(locked) > 128 {
				locked = make(map[int]*Machine)
			}

			freedCount := 0
			for index := len(r.InstanceList) - 1; index >= 0; index-- {
				instance := r.InstanceList[index]
				if instance == nil {
					continue
				}

				srcMachine := r.InstanceDeployedMachineMap[instance.InstanceId]
				machineId := status.InstanceDeployedMachineMap[instance.InstanceId].MachineId
				mappedMachineId := machineId
				if mapping[machineId] != 0 {
					mappedMachineId = mapping[machineId]
				}
				if mappedMachineId == srcMachine.MachineId {
					continue
				}

				if m.ConstraintCheck(instance) {
					r.CommandDeployInstance(instance, m)
					m.AddInstance(instance)
					freedCount++
					if freedCount > 128 {
						break
					}
				}
			}
		}

		if !hasSkipped {
			break
		}
	}

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
