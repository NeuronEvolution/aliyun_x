package bfs_v2

import (
	"fmt"
	"sort"
)

func (s *Strategy) debug() {
	sort.Slice(s.machineDeployList, func(i, j int) bool {
		if s.machineDeployList[i].LevelConfig.Cpu > s.machineDeployList[j].LevelConfig.Cpu {
			return true
		} else if s.machineDeployList[i].LevelConfig.Cpu == s.machineDeployList[j].LevelConfig.Cpu {
			return s.machineDeployList[i].GetCpuCost() > s.machineDeployList[j].GetCpuCost()
		} else {
			return false
		}
	})
	for _, m := range s.machineDeployList {
		fmt.Printf("cost=%f %f\n", m.GetCpuCost(), m.GetLinearCpuCost(m.LevelConfig.Cpu))
		if m.GetLinearCpuCost(m.LevelConfig.Cpu) < 0.9 {
			m.DebugPrint()
		} else {
			m.Resource.DebugPrint()
		}
	}
}
