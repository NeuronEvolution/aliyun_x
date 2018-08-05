package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

func (s *Strategy) mergeFinal() {
	startCost := s.R.CalculateTotalCostScore()
	fmt.Println("mergeFinal start cpu cost", startCost)

	currentCost := startCost
	loop := 0
	deadLoop := 0
	for ; ; loop++ {
		machines := make([]*cloud.Machine, 0)
		machinePool := append(make([]*cloud.Machine, 0), s.machineDeployList...)

		//CPU头部和尾部概率大
		cloud.SortMachineByCpuCost(machinePool)
		machinesByCpu := s.randMachinesBig2Big(machinePool, 2)
		machines = append(machines, machinesByCpu...)
		machinePool = cloud.MachinesRemove(machinePool, machinesByCpu)

		//CPU方差头部概率大
		sort.Slice(machinePool, func(i, j int) bool {
			return machinePool[i].GetCpuDerivation() > machinePool[j].GetCpuDerivation()
		})
		machinesByDerivation := s.randMachinesBig2Small(machinePool, 0)
		machines = append(machines, machinesByDerivation...)
		machinePool = cloud.MachinesRemove(machinePool, machinesByDerivation)

		ok, delta := s.BestMergeMachines(machines, deadLoop)
		if !ok {
			fmt.Println("mergeFinal dead loop", deadLoop)
			deadLoop++
			if deadLoop > 128 {
				break
			}

			continue
		}

		deadLoop = 0

		currentCost += delta

		fmt.Printf("mergeFinal loop %d %f %f\n", loop, startCost, currentCost)
	}

	fmt.Printf("mergeFinal end %d %f %f\n", loop, startCost, s.R.CalculateTotalCostScore())
}
