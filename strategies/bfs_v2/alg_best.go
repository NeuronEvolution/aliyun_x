package bfs_v2

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
	"sync"
)

func (s *Strategy) BatchBestMergeMachines(machines []*cloud.Machine, deadLoop int) (has bool, delta float64) {
	wg := &sync.WaitGroup{}
	max := len(machines)
	if len(machines)%2 == 1 {
		max = len(machines) - 1
	}
	for i := 0; i < max; i += 2 {
		batchMachines := []*cloud.Machine{machines[i], machines[i+1]}
		wg.Add(1)
		go func() {
			defer wg.Done()

			ok, d := s.BestMergeMachines(batchMachines, deadLoop)
			if ok {
				has = true
				bestLock2.Lock()
				defer bestLock2.Unlock()
				delta += d
			}
		}()
	}

	wg.Wait()

	return has, delta
}

var bestLock = &sync.Mutex{}
var bestLock2 = &sync.Mutex{}

func (s *Strategy) BestMergeMachines(machines []*cloud.Machine, deadLoop int) (has bool, delta float64) {
	instances := make([]*cloud.Instance, 0)
	for _, m := range machines {
		instances = append(instances, m.InstanceArray[:m.InstanceArrayCount]...)
	}

	cost := float64(0)
	for _, m := range machines {
		cost += m.GetCpuCost()
	}

	bestPos, bestCost := s.Best(machines, instances, deadLoop)
	if bestCost >= cost {
		return false, 0
	}

	bestLock.Lock()
	defer bestLock.Unlock()

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

func (s *Strategy) bestCheck(pos []int, machines []*cloud.Machine, instances []*cloud.Instance) {
	machineCount := len(machines)
	resources := make([]*cloud.Resource, machineCount)
	appCounts := make([]*cloud.AppCountCollection, machineCount)
	for i := 0; i < machineCount; i++ {
		resources[i] = &cloud.Resource{}
		appCounts[i] = cloud.NewAppCountCollection()
	}
	for i, instance := range instances {
		machineIndex := pos[i]
		m := machines[machineIndex]

		if !cloud.ConstraintCheckResourceLimit(resources[machineIndex], &instance.Config.Resource, m.LevelConfig, 1) ||
			!cloud.ConstraintCheckAppInterferenceAddInstance(instance.Config.AppId, appCounts[machineIndex], s.R.AppInterferenceConfigMap) {
			panic("111")
		}

		resources[machineIndex].AddResource(&instance.Config.Resource)
		appCounts[machineIndex].Add(instance.Config.AppId)
	}
}

func (s *Strategy) Best(machines []*cloud.Machine, instances []*cloud.Instance, deadLoop int) (bestPos []int, bestCost float64) {
	totalLoopLimit := 1024 * 32 * int(math.Pow(float64(2), float64(deadLoop)))
	if totalLoopLimit > 1024*1024 {
		totalLoopLimit = 1024 * 1024
	}

	//fmt.Printf("best machines=%d,instance=%d\n", len(machines), len(instances))

	//for _, instance := range instances {
	//	fmt.Println("instance", instance.InstanceId, instance.Config.AppId)
	//}

	machineCount := len(machines)
	instanceCount := len(instances)
	pos := make([]int, instanceCount)
	bestPos = make([]int, instanceCount)
	bestCost = math.MaxFloat64
	resources := make([]*cloud.Resource, machineCount)
	appCounts := make([]*cloud.AppCountCollection, machineCount)
	for i := 0; i < machineCount; i++ {
		resources[i] = &cloud.Resource{}
		appCounts[i] = cloud.NewAppCountCollection()
	}

	totalLoop := 0

	for instanceIndex := 0; instanceIndex < instanceCount; instanceIndex++ {
		//fmt.Println("POS", pos)
		//fmt.Println(appCounts[pos[instanceIndex]].List[:appCounts[pos[instanceIndex]].ListCount])
		instance := instances[instanceIndex]
		added := false
		for ; pos[instanceIndex] < machineCount; pos[instanceIndex]++ {
			totalLoop++

			machineIndex := pos[instanceIndex]
			m := machines[machineIndex]
			if !cloud.ConstraintCheckResourceLimit(resources[machineIndex], &instance.Config.Resource, m.LevelConfig, 1) ||
				!cloud.ConstraintCheckAppInterferenceAddInstance(instance.Config.AppId, appCounts[machineIndex], s.R.AppInterferenceConfigMap) {
				continue
			}

			//fmt.Println("ADD", instance.Config.AppId, appCounts[machineIndex].List[:appCounts[machineIndex].ListCount])
			resources[machineIndex].AddResource(&instance.Config.Resource)
			appCounts[machineIndex].Add(instance.Config.AppId)
			//fmt.Println("ADD", instance.Config.AppId, appCounts[machineIndex].List[:appCounts[machineIndex].ListCount])

			added = true
			break
		}

		if added {
			//有效解,回退
			if instanceIndex == instanceCount-1 {
				//s.bestCheck(pos, machines, instances)
				//fmt.Println("RESULT", pos)
				totalCost := float64(0)
				for machineIndex, r := range resources {
					totalCost += r.GetCpuCost(machines[machineIndex].LevelConfig.Cpu)
				}

				//最优解
				if totalCost < bestCost {
					//fmt.Println("BEST", bestCost, totalCost)
					bestCost = totalCost
					for i, v := range pos {
						bestPos[i] = v
					}
					//fmt.Println(bestPos)
				}

				//回退
				//fmt.Println("BACK")
				//fmt.Println(pos[instanceIndex], pos, instance.Config.AppId)
				//fmt.Println(appCounts[pos[instanceIndex]].List[:appCounts[pos[instanceIndex]].ListCount])
				resources[pos[instanceIndex]].RemoveResource(&instance.Config.Resource)
				appCounts[pos[instanceIndex]].Remove(instance.Config.AppId)
				//fmt.Println(appCounts[pos[instanceIndex]].List[:appCounts[pos[instanceIndex]].ListCount])
				pos[instanceIndex] = 0
			}
		} else {
			//回退
			pos[instanceIndex] = 0
		}

		end := false
		if !added || instanceIndex == instanceCount-1 {
			for {
				//已到最后
				instanceIndex--
				if instanceIndex < 0 {
					end = true
					break
				}

				//fmt.Println("INC")
				//fmt.Println(pos[instanceIndex], pos, instances[instanceIndex].Config.AppId)
				//fmt.Println(appCounts[pos[instanceIndex]].List[:appCounts[pos[instanceIndex]].ListCount])
				resources[pos[instanceIndex]].RemoveResource(&instances[instanceIndex].Config.Resource)
				appCounts[pos[instanceIndex]].Remove(instances[instanceIndex].Config.AppId)
				//fmt.Println(appCounts[pos[instanceIndex]].List[:appCounts[pos[instanceIndex]].ListCount])

				pos[instanceIndex]++
				if pos[instanceIndex] < machineCount {
					//进位成功
					instanceIndex--
					break
				} else {
					pos[instanceIndex] = 0
				}
			}
		}

		if end || (instanceCount > 20 && totalLoop > totalLoopLimit) {
			break
		}
	}

	//fmt.Println("BEST total loop", totalLoop)

	return bestPos, bestCost
}
