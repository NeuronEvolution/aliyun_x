package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
	"sort"
)

func (s *Strategy) mergeBestFit(machines []*cloud.Machine, instances []*cloud.Instance) (pos []int, cost float64) {
	machineCount := len(machines)
	instanceCount := len(instances)
	pos = make([]int, instanceCount)
	resources := make([]*cloud.Resource, machineCount)
	appCounts := make([]*cloud.AppCountCollection, machineCount)
	for i := 0; i < machineCount; i++ {
		resources[i] = &cloud.Resource{}
		appCounts[i] = cloud.NewAppCountCollection()
	}

	for i, instance := range instances {
		minCost := math.MaxFloat64
		minMachineIndex := -1
		for machineIndex := 0; machineIndex < machineCount; machineIndex++ {
			m := machines[machineIndex]
			if !cloud.ConstraintCheckResourceLimit(
				resources[machineIndex], &instance.Config.Resource, m.LevelConfig, 0.5) {
				continue
			}

			if !cloud.ConstraintCheckAppInterference(appCounts[machineIndex], s.R.AppInterferenceConfigMap) {
				continue
			}

			cost := resources[machineIndex].GetCostWithInstance(instance, m.LevelConfig.Cpu)
			if cost < minCost {
				//fmt.Println("cost", minCost, cost, machineIndex)
				minCost = cost
				minMachineIndex = machineIndex
				//fmt.Println("minMachineIndex", minMachineIndex)
			}
		}

		//fmt.Println("mergeBestFit A ", i, minCost, minMachineIndex)

		if minMachineIndex == -1 {
			for machineIndex := 0; machineIndex < machineCount; machineIndex++ {
				m := machines[machineIndex]
				if !cloud.ConstraintCheckResourceLimit(
					resources[machineIndex], &instance.Config.Resource, m.LevelConfig, 1) {
					continue
				}

				if !cloud.ConstraintCheckAppInterference(appCounts[machineIndex], s.R.AppInterferenceConfigMap) {
					continue
				}

				cost := resources[machineIndex].GetCostWithInstance(instance, m.LevelConfig.Cpu)
				if cost < minCost {
					//fmt.Println("mergeBestFit B ", i, minCost, minMachineIndex)
					minCost = cost
					minMachineIndex = machineIndex
				}
			}
		}

		if minMachineIndex == -1 {
			return nil, math.MaxFloat64
		}

		resources[minMachineIndex].AddResource(&instance.Config.Resource)
		appCounts[minMachineIndex].Add(instance.Config.AppId)
		pos[i] = minMachineIndex

		//fmt.Println("mergeBestFit C ", i, minCost, minMachineIndex, pos)
	}

	for i, r := range resources {
		cost += r.GetCpuCost(machines[i].LevelConfig.Cpu)
	}

	return pos, cost
}

func (s *Strategy) mergeBestFitMachines(machines []*cloud.Machine) (has bool, delta float64) {
	instances := make([]*cloud.Instance, 0)
	for _, m := range machines {
		instances = append(instances, m.InstanceArray[:m.InstanceArrayCount]...)
	}

	cost := float64(0)
	for _, m := range machines {
		cost += m.GetCpuCost()
	}

	sort.Slice(machines, func(i, j int) bool {
		return machines[i].LevelConfig.Cpu > machines[j].LevelConfig.Cpu
	})
	cloud.SortInstanceByTotalMaxLowWithInference(instances, 4)
	bestPos, bestCost := s.mergeBestFit(machines, instances)
	fmt.Println("mergeBestFitMachines", bestCost, cost)
	if bestCost >= cost {
		return false, 0
	}

	//将所有实例迁出
	for _, m := range machines {
		for _, inst := range cloud.InstancesCopy(m.InstanceArray[:m.InstanceArrayCount]) {
			m.RemoveInstance(inst.InstanceId)
		}
	}

	for i, instance := range instances {
		m := machines[bestPos[i]]
		if !m.ConstraintCheck(instance, m.LevelConfig.Cpu) {
			panic("ConstraintCheck")
		}
		m.AddInstance(instance)
	}

	return true, bestCost - cost
}

func (s *Strategy) merge() {
	startCost := s.R.CalculateTotalCostScore()
	fmt.Println("merge start cpu cost", startCost)

	currentCost := startCost
	loop := 0
	for ; loop < 1000; loop++ {
		cloud.SortMachineByCpuCost(s.machineDeployList)
		machinesByCpu := s.randMachinesBig2Big(s.machineDeployList, 128)
		has, delta := s.mergeBestFitMachines(machinesByCpu)
		if !has {
			fmt.Printf("merge loop failed %d %f\n", loop, startCost)
			continue
		}

		currentCost += delta

		fmt.Printf("merge loop %d %f %f\n", loop, startCost, currentCost)
	}

	fmt.Printf("merge end %d %f %f\n", loop, startCost, s.R.CalculateTotalCostScore())
}
