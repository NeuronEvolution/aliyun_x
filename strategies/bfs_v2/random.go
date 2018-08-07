package bfs_v2

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
)

func (s *Strategy) randBig2Big(count int) []float64 {
	pTable := make([]float64, count)
	for i := 0; i < count; i++ {
		pTable[i] = math.Exp((math.Abs(float64(i)-float64(count)/2) - float64(count)/2) /
			(float64(count) / 8))
		if i > 0 {
			pTable[i] += pTable[i-1]
		}
	}

	return pTable
}

func (s *Strategy) randMachinesBig2Big(pool []*cloud.Machine, pTable []float64, count int) (machines []*cloud.Machine) {
	machines = make([]*cloud.Machine, 0)
	maxP := pTable[len(pTable)-1]
	for i := 0; i < count; i++ {
		r := s.R.Rand.Float64() * maxP
		for machineIndex, p := range pTable {
			if p < r {
				continue
			}

			if cloud.MachinesContains(machines, pool[machineIndex].MachineId) {
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

func (s *Strategy) randMachinesBig2Small(pool []*cloud.Machine, count int) (machines []*cloud.Machine) {
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

			if cloud.MachinesContains(machines, pool[machineIndex].MachineId) {
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
