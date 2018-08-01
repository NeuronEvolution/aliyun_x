package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
	"math/rand"
	"sort"
)

//TODO 机器的概率分布
//TODO 实例的概率分布
//TODO 一次调整多个实例

const SALoopCount = 100000

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

type SAContext struct {
	Machines  []*cloud.Machine
	Instances []*cloud.Instance

	Rand            *rand.Rand
	CurrentMachines []*SAMachine
	CurrentMap      map[*cloud.Instance]*SAMachine
	BestCost        float64
	BestMap         map[*cloud.Instance]*SAMachine
	InferenceMap    [][cloud.MaxAppId]int
}

func NewSAContext(machines []*cloud.Machine, inferenceMap [][cloud.MaxAppId]int) *SAContext {
	ctx := &SAContext{}
	ctx.Machines = machines
	ctx.InferenceMap = inferenceMap

	return ctx
}

func (ctx *SAContext) init() {
	machineCount := len(ctx.Machines)
	ctx.Instances = make([]*cloud.Instance, 0)
	for _, m := range ctx.Machines {
		ctx.Instances = append(ctx.Instances, m.InstanceArray[:m.InstanceArrayCount]...)
	}
	ctx.Rand = rand.New(rand.NewSource(0))
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

func (ctx *SAContext) newStatus() (instance *cloud.Instance, oldMachine *SAMachine, newMachine *SAMachine) {
	for i := 0; i < 100; i++ {
		instanceRand := ctx.Rand.Intn(len(ctx.Instances))
		instance := ctx.Instances[instanceRand]
		oldMachine = ctx.CurrentMap[instance]

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
			return instance, oldMachine, newMachine
		}
	}

	return nil, nil, nil
}

func (ctx *SAContext) Run() {
	ctx.init()

	T := float64(1000000)
	r := float64(0.999)
	for i := 0; i < SALoopCount; i++ {
		fmt.Printf("SA loop %d %f\n", i, T)
		instance, oldMachine, newMachine := ctx.newStatus()
		if instance == nil {
			fmt.Println("SA new state failed")
			break
		}

		oldCost := oldMachine.GetCpuCost() + newMachine.GetCpuCost()
		oldMachine.Remove(instance)
		newMachine.Add(instance)
		newCost := oldMachine.GetCpuCost() + newMachine.GetCpuCost()

		accept := false
		dE := oldCost - newCost
		if dE > 0 {
			//fmt.Println("SA accept", dE)
			accept = true
		} else {
			if math.Exp(dE/T) > ctx.Rand.Float64() {
				accept = true
			}
		}

		if accept {
			ctx.CurrentMap[instance] = newMachine
			cost := float64(0)
			for _, m := range ctx.CurrentMachines {
				cost += m.GetCpuCost()
			}
			if cost < ctx.BestCost {
				fmt.Println("SA best", ctx.BestCost, cost)
				ctx.BestCost = cost
				ctx.BestMap = make(map[*cloud.Instance]*SAMachine)
				for instance, m := range ctx.CurrentMap {
					ctx.BestMap[instance] = m
				}
			}
		} else {
			newMachine.Remove(instance)
			oldMachine.Add(instance)
		}

		T = r * T
	}
}

func (s *Strategy) mergeMachineSA(machines []*cloud.Machine) bool {
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

	ctx := NewSAContext(machines, s.R.AppInterferenceConfigMap)
	ctx.Run()
	if ctx.BestCost >= cost {
		fmt.Printf("mergeMachineSA failed,cost=%f,sa cost=%f\n", cost, ctx.BestCost)
		return false
	}

	fmt.Printf("mergeMachineSA cost=%f best=%f", cost, ctx.BestCost)

	//将所有实例迁出
	for _, m := range machines {
		for _, inst := range cloud.InstancesCopy(m.InstanceArray[:m.InstanceArrayCount]) {
			m.RemoveInstance(inst.InstanceId)
		}
	}

	cloud.SetDebug(true)
	for instance, mSA := range ctx.BestMap {
		//mSA.Machine.DebugPrint()
		if !mSA.Machine.ConstraintCheck(instance, 1) {
			panic("mergeMachineSA ConstraintCheck failed")
		}

		mSA.Machine.AddInstance(instance)
	}
	cloud.SetDebug(false)

	return true
}
