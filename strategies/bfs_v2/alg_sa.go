package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
	"math/rand"
	"sort"
)

//TODO 实例的概率分布
//TODO 分批优化

const SALoopCount = 1000000
const SATemperature = 100000
const SARatio = 0.9995
const SANewStatusMoveCount = 4
const SANewStatusRetry = 100

type SAMachine struct {
	Machine       *cloud.Machine
	ResourceLimit *cloud.MachineLevelConfig
	Resource      *cloud.Resource
	AppCount      *cloud.AppCountCollection
}

func NewSAMachine(machine *cloud.Machine) *SAMachine {
	m := &SAMachine{}
	m.Machine = machine
	m.ResourceLimit = machine.LevelConfig
	m.Resource = &cloud.Resource{}
	m.AppCount = cloud.NewAppCountCollection()

	return m
}

func (m *SAMachine) Add(instance *cloud.Instance) {
	m.Resource.AddResource(&instance.Config.Resource)
	m.AppCount.Add(instance.Config.AppId)
}

func (m *SAMachine) Remove(instance *cloud.Instance) {
	m.Resource.RemoveResource(&instance.Config.Resource)
	m.AppCount.Remove(instance.Config.AppId)
}

func (m *SAMachine) GetCpuCost() float64 {
	return m.Resource.GetCpuCost(m.ResourceLimit.Cpu)
}

func (m *SAMachine) ConstraintCheck(instance *cloud.Instance, inferenceMap [][cloud.MaxAppId]int) bool {
	if !cloud.ConstraintCheckResourceLimit(m.Resource, &instance.Config.Resource, m.ResourceLimit, 1) {
		return false
	}

	if !cloud.ConstraintCheckAppInterferenceAddInstance(instance.Config.AppId, m.AppCount, inferenceMap) {
		return false
	}

	return true
}

func (m *SAMachine) Reset() {
	m.Resource.Reset()
	m.AppCount.Reset()
}

func (m *SAMachine) Debug() {
	m.Resource.DebugPrint()
	m.AppCount.Debug()
}

func (m *SAMachine) Validate() {

}

type SAMove struct {
	Instance   *cloud.Instance
	OldMachine *SAMachine
	NewMachine *SAMachine
}

type SAContext struct {
	Machines  []*cloud.Machine
	Instances []*cloud.Instance

	Rand            *rand.Rand
	CurrentMachines []*SAMachine
	CurrentMap      map[*cloud.Instance]*SAMachine
	BestCost        float64
	BestMap         map[*cloud.Instance]*SAMachine
	InferenceMap    [][cloud.MaxAppId]int
	HasBest         bool
}

func NewSAContext(rand *rand.Rand, machines []*cloud.Machine, inferenceMap [][cloud.MaxAppId]int) *SAContext {
	ctx := &SAContext{}
	ctx.Rand = rand
	ctx.Machines = machines
	ctx.InferenceMap = inferenceMap

	return ctx
}

func (ctx *SAContext) init() {
	ctx.HasBest = false
	machineCount := len(ctx.Machines)
	ctx.Instances = make([]*cloud.Instance, 0)
	for _, m := range ctx.Machines {
		ctx.Instances = append(ctx.Instances, m.InstanceArray[:m.InstanceArrayCount]...)
	}
	ctx.CurrentMachines = make([]*SAMachine, machineCount)
	ctx.CurrentMap = make(map[*cloud.Instance]*SAMachine)
	ctx.BestMap = make(map[*cloud.Instance]*SAMachine)
	for i, m := range ctx.Machines {
		mSA := NewSAMachine(m)
		for _, instance := range m.InstanceArray[:m.InstanceArrayCount] {
			mSA.Add(instance)
			ctx.CurrentMap[instance] = mSA
			ctx.BestMap[instance] = mSA
		}
		ctx.CurrentMachines[i] = mSA
		ctx.BestCost += mSA.GetCpuCost()
	}
}

func (ctx *SAContext) validate() {
	for _, m := range ctx.CurrentMachines {
		m.Validate()
	}

	machineResourceMap := make(map[*SAMachine]*SAMachine)
	for instance, machine := range ctx.CurrentMap {
		resource := machineResourceMap[machine]
		if resource == nil {
			resource = NewSAMachine(machine.Machine)
			machineResourceMap[machine] = resource
		}

		if !resource.ConstraintCheck(instance, ctx.InferenceMap) {
			panic("ConstraintCheck")
		}

		resource.Add(instance)
	}

	machineResourceMap = make(map[*SAMachine]*SAMachine)
	for instance, machine := range ctx.BestMap {
		resource := machineResourceMap[machine]
		if resource == nil {
			resource = NewSAMachine(machine.Machine)
			machineResourceMap[machine] = resource
		}

		if !resource.ConstraintCheck(instance, ctx.InferenceMap) {
			panic("ConstraintCheck")
		}

		resource.Add(instance)
	}
}

func (ctx *SAContext) newStatus() (moves []*SAMove, delta float64) {
	for i := 0; i < SANewStatusMoveCount; i++ {
		has := false
		for loop := 0; loop < SANewStatusRetry; loop++ {
			instanceRand := ctx.Rand.Intn(len(ctx.Instances))
			instance := ctx.Instances[instanceRand]

			already := false
			for _, move := range moves {
				if move.Instance == instance {
					already = true
				}
			}
			if already {
				continue
			}

			var newMachine *SAMachine
			oldMachine := ctx.CurrentMap[instance]
			machineRand := ctx.Rand.Intn(len(ctx.Machines))
			pos := machineRand
			for {
				m := ctx.CurrentMachines[pos]
				if m != oldMachine && m.ConstraintCheck(instance, ctx.InferenceMap) {
					newMachine = m
					break
				}

				pos++
				if pos == len(ctx.Machines) {
					pos = 0
				}
				if pos == machineRand {
					break
				}
			}

			if newMachine != nil {
				oldCost := oldMachine.GetCpuCost() + newMachine.GetCpuCost()
				oldMachine.Remove(instance)
				newMachine.Add(instance)
				ctx.CurrentMap[instance] = newMachine
				newCost := oldMachine.GetCpuCost() + newMachine.GetCpuCost()
				delta += oldCost - newCost
				moves = append(moves, &SAMove{Instance: instance, OldMachine: oldMachine, NewMachine: newMachine})
				has = true
				break
			}
		}
		if !has {
			break
		}
	}

	return moves, delta
}

func (ctx *SAContext) Run() {
	ctx.init()

	T := float64(SATemperature)
	r := float64(SARatio)
	for loop := 0; loop < SALoopCount; loop++ {
		if loop > 0 && loop%100000 == 0 {
			fmt.Printf("SA loop %d %f\n", loop, T)
		}

		moves, dE := ctx.newStatus()
		if len(moves) == 0 {
			fmt.Println("SA new state failed")
			break
		}

		accept := false
		if dE > 0 {
			//fmt.Println("SA accept", dE)
			accept = true
		} else {
			if math.Exp(dE/T) > ctx.Rand.Float64() {
				accept = true
			}
		}

		if accept {
			cost := float64(0)
			for _, m := range ctx.CurrentMachines {
				cost += m.GetCpuCost()
			}
			if cost < ctx.BestCost {
				fmt.Println("SA best", ctx.BestCost, cost)
				ctx.HasBest = true
				ctx.BestCost = cost
				ctx.BestMap = make(map[*cloud.Instance]*SAMachine)
				for instance, m := range ctx.CurrentMap {
					ctx.BestMap[instance] = m
				}
			}
		} else {
			for i := len(moves) - 1; i >= 0; i-- {
				move := moves[i]
				move.NewMachine.Remove(move.Instance)
				move.OldMachine.Add(move.Instance)
				ctx.CurrentMap[move.Instance] = move.OldMachine
			}
		}

		T = r * T
	}
}

func (s *Strategy) mergeMachineSA(machines []*cloud.Machine) (ok bool, delta float64) {
	//for _, m := range machines {
	//	m.DebugPrint()
	//}

	cost := float64(0)
	for _, m := range machines {
		cost += m.GetCpuCost()
	}

	sort.Slice(machines, func(i, j int) bool {
		return machines[i].Disk > machines[j].Disk
	})

	//纪录当前状态
	instanceMachineMap := make(map[*cloud.Instance]*cloud.Machine)
	for _, m := range machines {
		for _, instance := range cloud.InstancesCopy(m.InstanceArray[:m.InstanceArrayCount]) {
			instanceMachineMap[instance] = m
		}
	}

	ctx := NewSAContext(s.R.Rand, machines, s.R.AppInterferenceConfigMap)
	ctx.Run()
	if !ctx.HasBest {
		return false, 0
	}

	fmt.Printf("mergeMachineSA cost=%f best=%f\n", cost, ctx.BestCost)

	//将所有实例迁出
	for _, m := range machines {
		for _, inst := range cloud.InstancesCopy(m.InstanceArray[:m.InstanceArrayCount]) {
			m.RemoveInstance(inst.InstanceId)
		}
	}

	for instance, mSA := range ctx.BestMap {
		mSA.Machine.AddInstance(instance)
	}

	return true, ctx.BestCost - cost
}
