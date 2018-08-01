package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

func (s *Strategy) merge() {
	fmt.Println("merge start cpu cost", s.R.CalculateTotalCostScore())

	sort.Slice(s.machineDeployList, func(i, j int) bool {
		m1 := s.machineDeployList[i]
		m2 := s.machineDeployList[j]
		cpu1 := m1.GetCpuCostReal()
		cpu2 := m2.GetCpuCostReal()
		linearCpu1 := m1.GetLinearCpuCost(m1.LevelConfig.Cpu)
		linearCpu2 := m2.GetLinearCpuCost(m2.LevelConfig.Cpu)
		if cpu1 > 1.01 || cpu2 > 1.01 {
			return cpu1 > cpu2
		}

		return linearCpu1 > linearCpu2
	})

	machines := make([]*cloud.Machine, 0)
	machines = append(s.machineDeployList[:10], s.machineDeployList[len(s.machineDeployList)-10:]...)

	s.mergeMachineSA(machines)
}
