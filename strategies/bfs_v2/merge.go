package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

func (s *Strategy) merge() {
	startCost := s.R.CalculateTotalCostScore()
	fmt.Println("merge start cpu cost", startCost)

	currentCost := startCost
	count := 32
	loop := 0
	deadLoop := 0
	for ; ; loop++ {
		machinePool := append(make([]*cloud.Machine, 0), s.machineDeployList...)

		//CPU头部和尾部
		sort.Slice(machinePool, func(i, j int) bool {
			m1 := machinePool[i]
			m2 := machinePool[j]
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
		machines = append(machines, machinePool[:count]...)
		machines = append(machines, machinePool[len(machinePool)-count:]...)
		machinePool = machinePool[count : len(machinePool)-count]
		//fmt.Println("machinePool len 1", len(machinePool))

		//CPU方差头部
		sort.Slice(machinePool, func(i, j int) bool {
			return machinePool[i].GetCpuDerivation() > machinePool[j].GetCpuDerivation()
		})
		machines = append(machines, machinePool[:count]...)
		machinePool = machinePool[count:]
		//fmt.Println("machinePool len 2", len(machinePool))

		//随机取若干台
		for i := 0; i < count; i++ {
			//fmt.Println("machinePool len 3", len(machinePool))
			m := machinePool[s.R.Rand.Intn(len(machinePool))]
			machines = append(machines, m)
			machinePool = cloud.MachinesRemove(machinePool, []*cloud.Machine{m})
		}

		ok, delta := s.mergeMachineSA(machines)
		if !ok {
			fmt.Println("merge dead loop", deadLoop)
			deadLoop++
			if deadLoop > 10 {
				break
			}

			continue
		}

		deadLoop = 0

		currentCost += delta

		fmt.Printf("merge loop %d %f %f\n", loop, startCost, currentCost)
	}

	fmt.Printf("merge end %d %f %f\n", loop, startCost, s.R.CalculateTotalCostScore())
}
