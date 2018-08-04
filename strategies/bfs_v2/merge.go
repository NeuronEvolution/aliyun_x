package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
	"sort"
)

func (s *Strategy) randMachinesByCpu(pool []*cloud.Machine, count int) (machines []*cloud.Machine) {
	cloud.SortMachineByCpuCost(pool)

	machines = make([]*cloud.Machine, 0)
	machineCount := len(pool)
	pTable := make([]float64, machineCount)
	for i := 0; i < machineCount; i++ {
		pTable[i] = math.Exp((math.Abs(float64(i)-float64(machineCount)/2) - float64(machineCount)/2) /
			(float64(machineCount) / 8))
		if i > 0 {
			pTable[i] += pTable[i-1]
		}
	}
	maxP := pTable[len(pTable)-1]
	for i := 0; i < count; i++ {
		r := s.R.Rand.Float64() * maxP
		for machineIndex, p := range pTable {
			if p < r {
				continue
			}

			if cloud.MachinesContaines(machines, pool[machineIndex].MachineId) {
				if i == count-1 {
					i = -1
				}
				continue
			}

			machines = append(machines, pool[machineIndex])
			break
		}
	}

	return machines
}

func (s *Strategy) randMachinesByDerivation(pool []*cloud.Machine, count int) (machines []*cloud.Machine) {
	sort.Slice(pool, func(i, j int) bool {
		return pool[i].GetCpuDerivation() > pool[j].GetCpuDerivation()
	})

	machines = make([]*cloud.Machine, 0)
	machineCount := len(pool)
	pTable := make([]float64, machineCount)
	for i := 0; i < machineCount; i++ {
		pTable[i] = math.Exp(-float64(i) * 8 / float64(machineCount))
		if i > 0 {
			pTable[i] += pTable[i-1]
		}
	}
	maxP := pTable[len(pTable)-1]
	for i := 0; i < count; i++ {
		r := s.R.Rand.Float64() * maxP
		for machineIndex, p := range pTable {
			if p < r {
				continue
			}

			if cloud.MachinesContaines(machines, pool[machineIndex].MachineId) {
				if i == count-1 {
					i = -1
				}
				continue
			}

			machines = append(machines, pool[machineIndex])
			break
		}
	}

	return machines
}

func (s *Strategy) merge() {
	startCost := s.R.CalculateTotalCostScore()
	fmt.Println("merge start cpu cost", startCost)

	currentCost := startCost
	loop := 0
	deadLoop := 0
	for ; ; loop++ {
		machines := make([]*cloud.Machine, 0)
		machinePool := append(make([]*cloud.Machine, 0), s.machineDeployList...)

		//CPU头部和尾部概率大
		machinesByCpu := s.randMachinesByCpu(machinePool, 2)
		machines = append(machines, machinesByCpu...)
		machinePool = cloud.MachinesRemove(machinePool, machinesByCpu)

		//CPU方差头部概率大
		machinesByDerivation := s.randMachinesByDerivation(machinePool, 0)
		machines = append(machines, machinesByDerivation...)
		machinePool = cloud.MachinesRemove(machinePool, machinesByDerivation)

		ok, delta := s.BestMergeMachines(machines)
		if !ok {
			fmt.Println("merge dead loop", deadLoop)
			deadLoop++
			if deadLoop > 32 {
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
