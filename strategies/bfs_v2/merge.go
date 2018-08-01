package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

func (s *Strategy) mergeMachine(machines []*cloud.Machine) bool {
	cost := float64(0)
	instances := make([]*cloud.Instance, 0)
	for _, m := range machines {
		//m.DebugPrint()
		cost += m.GetCpuCost()
		instances = append(instances, m.InstanceArray[:m.InstanceArrayCount]...)
	}

	sort.Slice(machines, func(i, j int) bool {
		return machines[i].Disk > machines[j].Disk
	})

	//PSO优化
	ctx := &PsoContext{Machines: machines, Instances: instances, InferenceMap: s.R.AppInterferenceConfigMap}
	ctx.Run()
	if ctx.BestCost >= cost {
		fmt.Printf("mergeMachine failed,cost=%f best=%f\n", cost, ctx.BestCost)
		//return false
	}

	fmt.Printf("mergeMachine ok,cost=%f best=%f\n", cost, ctx.BestCost)

	//纪录当前状态
	instanceMachineMap := make(map[*cloud.Instance]*cloud.Machine)
	for _, m := range machines {
		for _, instance := range cloud.InstancesCopy(m.InstanceArray[:m.InstanceArrayCount]) {
			instanceMachineMap[instance] = m
		}
	}

	//将所有实例迁出
	for _, m := range machines {
		for _, inst := range cloud.InstancesCopy(m.InstanceArray[:m.InstanceArrayCount]) {
			m.RemoveInstance(inst.InstanceId)
		}
	}

	//使用PSO最优结果
	failed := false
	cloud.SetDebug(true)
	for instanceIndex, machineIndex := range ctx.BestParticle.BestPosition {
		m := machines[machineIndex]
		instance := instances[instanceIndex]
		if !m.ConstraintCheck(instance, 1) {
			failed = true
			fmt.Println("mergeMachine pso ConstraintCheck failed")
			m.DebugPrint()
			instance.Config.DebugPrint()
			break
		}
		m.AddInstance(instance)
	}
	cloud.SetDebug(false)

	//PSO最优结果无效，恢复到原状态
	if failed {
		for _, m := range machines {
			for _, inst := range cloud.InstancesCopy(m.InstanceArray[:m.InstanceArrayCount]) {
				m.RemoveInstance(inst.InstanceId)
			}
		}

		for instance, m := range instanceMachineMap {
			m.AddInstance(instance)
		}

		return false
	}

	return true
}

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

	for _, m := range machines {
		m.DebugPrint()
	}

	s.mergeMachine(machines)
}
