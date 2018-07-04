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

	for i := 0; ; i++ {
		fmt.Printf("ResourceManagement.MergeTo loop %d\n", i)
		hasSkippedInstance := false
		for _, instance := range r.InstanceList {
			if instance == nil {
				continue
			}

			srcMachine := r.InstanceDeployedMachineMap[instance.InstanceId]
			machineId := status.InstanceDeployedMachineMap[instance.InstanceId].MachineId
			mappedMachineId := machineId
			if mapping[machineId] != 0 {
				mappedMachineId = mapping[machineId]
			}
			if srcMachine.MachineId != mappedMachineId {
				targetMachine := r.MachineMap[mappedMachineId]
				if targetMachine.ConstraintCheck(instance) {
					r.CommandDeployInstance(instance, targetMachine)
				} else {
					hasSkippedInstance = true
				}
			} else {
				fmt.Printf("mapped %d\n", machineId)
			}
		}

		if !hasSkippedInstance {
			break
		}

		if i > 10 {
			return fmt.Errorf("ResourceManagement.MergeTo deadLoop\n")
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
