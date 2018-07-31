package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

func (s *Strategy) mergeMachine(machines []*cloud.Machine) bool {
	for _, m := range machines {
		m.DebugPrint()
	}

	//将所有实例迁出
	instances := make([]*cloud.Instance, 0)
	for _, m := range machines {
		instances = append(instances, m.InstanceArray[:m.InstanceArrayCount]...)
	}
	for _, m := range machines {
		for _, inst := range cloud.InstancesCopy(m.InstanceArray[:m.InstanceArrayCount]) {
			m.RemoveInstance(inst.InstanceId)
		}
	}

	sort.Slice(machines, func(i, j int) bool {
		return machines[i].Disk > machines[j].Disk
	})

	fmt.Println("instances ", len(instances))

	return false
}

func (s *Strategy) merge() {
	fmt.Println("merge start cpu cost", s.R.CalculateTotalCostScore())

	sort.Slice(s.machineDeployList, func(i, j int) bool {
		m1 := s.machineDeployList[i]
		m2 := s.machineDeployList[j]
		cpu1 := m1.GetCostReal()
		cpu2 := m2.GetCostReal()
		linearCpu1 := m1.GetLinearCpuCost(m1.LevelConfig.Cpu)
		linearCpu2 := m2.GetLinearCpuCost(m2.LevelConfig.Cpu)
		if cpu1 > 1.01 || cpu2 > 1.01 {
			return cpu1 > cpu2
		}

		return linearCpu1 > linearCpu2
	})

	s.mergeMachine([]*cloud.Machine{s.machineDeployList[0], s.machineDeployList[len(s.machineDeployList)-1]})
}
