package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
	"math/rand"
	"sort"
)

const PsoParticleCount = 256
const PsoLoopCount = 256

type Particle struct {
	Position     []int              //最优状态。第几个Instance在第几个Machine
	Machines     []*ParticleMachine //机器部署计算状态
	BestCost     float64            //最优得分
	BestPosition []int              //最优状态。第几个Instance在第几个Machine
}

func NewParticle() *Particle {
	p := &Particle{}

	return p
}

func (p *Particle) GetCost(inferenceMap [][cloud.MaxAppId]int) float64 {
	totalCost := float64(0)
	for _, m := range p.Machines {
		totalCost += m.GetCost(inferenceMap)
	}

	return totalCost
}

func (p *Particle) RandReset(ctx *PsoContext) {
	for _, m := range p.Machines {
		m.Reset()
	}

	p.BestCost = math.MaxFloat64
	for instanceIndex := 0; instanceIndex < len(ctx.Instances); instanceIndex++ {
		machineIndex := ctx.Rand.Intn(len(p.Machines))
		p.Position[instanceIndex] = machineIndex
		p.BestPosition[instanceIndex] = machineIndex
		p.Machines[machineIndex].Add(ctx.Instances[instanceIndex])
	}
}

type ParticleMachine struct {
	ResourceLimit *cloud.MachineLevelConfig
	Resource      *cloud.Resource
	AppCount      *cloud.AppCountCollection
}

func NewParticleMachine(resourceLimit *cloud.MachineLevelConfig) *ParticleMachine {
	m := &ParticleMachine{}
	m.ResourceLimit = resourceLimit
	m.Resource = &cloud.Resource{}
	m.AppCount = cloud.NewAppCountCollection()

	return m
}

func (m *ParticleMachine) Add(instance *cloud.Instance) {
	m.Resource.AddResource(&instance.Config.Resource)
	m.AppCount.Add(instance.Config.AppId)
}

func (m *ParticleMachine) Remove(instance *cloud.Instance) {
	m.Resource.RemoveResource(&instance.Config.Resource)
	m.AppCount.Remove(instance.Config.AppId)
}

func (m *ParticleMachine) Reset() {
	m.Resource.Reset()
	m.AppCount.Reset()
}

func (m *ParticleMachine) Debug() {
	m.Resource.DebugPrint()
	m.AppCount.Debug()
}

func (m *ParticleMachine) GetInference(inferenceMap [][cloud.MaxAppId]int) int {
	total := 0
	for _, v1 := range m.AppCount.List[:m.AppCount.ListCount] {
		for _, v2 := range m.AppCount.List[:m.AppCount.ListCount] {
			maxCount := inferenceMap[v1.AppId][v2.AppId]
			if maxCount == -1 {
				continue
			}

			if v1.AppId == v2.AppId {
				maxCount++
			}

			if v2.Count > maxCount {
				total += v2.Count - maxCount
			}
		}
	}

	return total * 100
}

func (m *ParticleMachine) GetCost(inferenceMap [][cloud.MaxAppId]int) float64 {
	cpuCost := m.Resource.GetCpuCost(m.ResourceLimit.Cpu)
	badCpu := false
	for _, c := range m.Resource.Cpu {
		if c > m.ResourceLimit.Mem+cloud.ConstraintE {
			cpuCost += c / m.ResourceLimit.Mem
			badCpu = true
		}
	}
	if badCpu {
		cpuCost += 1
	}

	inference := m.GetInference(inferenceMap)

	disk := float64(0)
	if m.Resource.Disk > m.ResourceLimit.Disk {
		disk += 1 + float64(m.Resource.Disk)/float64(m.ResourceLimit.Disk)
	}

	memCost := float64(0)
	badMem := false
	for _, mem := range m.Resource.Mem {
		if mem > m.ResourceLimit.Mem+cloud.ConstraintE {
			memCost += mem / m.ResourceLimit.Mem
			badMem = true
		}
	}
	if badMem {
		memCost += 1
	}

	p := 0
	if m.Resource.P > m.ResourceLimit.P {
		p += m.Resource.P - m.ResourceLimit.P
	}

	mCost := 0
	if m.Resource.M > m.ResourceLimit.M {
		mCost += m.Resource.M - m.ResourceLimit.M
	}

	pm := 0
	if m.Resource.PM > m.ResourceLimit.PM {
		pm += m.Resource.PM - m.ResourceLimit.PM
	}

	return cpuCost + disk + memCost + float64(inference) + float64(p) + float64(mCost) + float64(pm)
}

type PsoContext struct {
	Machines  []*cloud.Machine
	Instances []*cloud.Instance

	Rand         *rand.Rand
	Particles    []*Particle
	BestParticle *Particle
	BestCost     float64
	InferenceMap [][cloud.MaxAppId]int
}

func (ctx *PsoContext) init() {
	machineCount := len(ctx.Machines)
	instanceCount := len(ctx.Instances)

	ctx.Rand = rand.New(rand.NewSource(0))
	ctx.BestCost = math.MaxFloat64

	ctx.Particles = make([]*Particle, PsoParticleCount)
	for particleIndex := 0; particleIndex < len(ctx.Particles); particleIndex++ {
		particle := NewParticle()
		particle.Machines = make([]*ParticleMachine, machineCount)
		for machineIndex := 0; machineIndex < len(particle.Machines); machineIndex++ {
			particle.Machines[machineIndex] = NewParticleMachine(ctx.Machines[machineIndex].LevelConfig)
		}
		particle.Position = make([]int, instanceCount)
		particle.BestPosition = make([]int, instanceCount)
		particle.RandReset(ctx)
		ctx.Particles[particleIndex] = particle
		particle.BestCost = particle.GetCost(ctx.InferenceMap)
		if particle.BestCost < ctx.BestCost {
			ctx.BestCost = particle.BestCost
			ctx.BestParticle = particle
		}
	}

	fmt.Println("pso init best position", ctx.BestParticle.BestPosition)
	fmt.Println("pso init best cost ", ctx.BestCost)
}

func (ctx *PsoContext) distance(p1 []int, p2 []int) int {
	d := 0
	for i, v1 := range p1 {
		if p2[i] != v1 {
			d++
		}
	}

	return d
}

func (ctx *PsoContext) Run() {
	machineCount := len(ctx.Machines)
	instanceCount := len(ctx.Instances)
	fmt.Printf("PSO machine %d,instance %d\n", machineCount, instanceCount)

	ctx.init()

	bestPosition := ctx.BestParticle.BestPosition
	for loop := 0; loop < PsoLoopCount; loop++ {
		if loop > 0 && loop%100 == 0 {
			fmt.Println("PSO loop", loop)
		}

		for _, particle := range ctx.Particles {
			//更新粒子速度
			velocity := ctx.Rand.Float64()
			//fmt.Println("PSO velocity", velocity)

			//更新粒子位置
			for instanceIndex := 0; instanceIndex < instanceCount; instanceIndex++ {
				machineIndex := particle.Position[instanceIndex]
				newMachineIndex := bestPosition[instanceIndex]
				if machineIndex == newMachineIndex {
					continue
				}

				r := ctx.Rand.Float64()
				if r > velocity {
					continue
				}

				instance := ctx.Instances[instanceIndex]
				particle.Machines[machineIndex].Remove(instance)
				particle.Machines[newMachineIndex].Add(instance)
				particle.Position[instanceIndex] = newMachineIndex
				break
			}

			//粒子和最优状态重合，重新生成新粒子
			if ctx.distance(particle.Position, bestPosition) <= 0 {
				if particle != ctx.BestParticle {
					particle.RandReset(ctx)
				}
			}

			//更新粒子最优状态
			particleCost := particle.GetCost(ctx.InferenceMap)
			if particleCost < particle.BestCost {
				//fmt.Printf("PSO update p best %f %f\n", particle.BestCost, particleCost)
				particle.BestCost = particleCost
				for instanceIndex := 0; instanceIndex < instanceCount; instanceIndex++ {
					particle.BestPosition[instanceIndex] = particle.Position[instanceIndex]
				}

				//更新全局最优状态
				if particle.BestCost < ctx.BestCost {
					fmt.Printf("PSO update g best %f %f\n", ctx.BestCost, particle.BestCost)
					ctx.BestCost = particle.BestCost
					ctx.BestParticle = particle
					bestPosition = ctx.BestParticle.BestPosition
				}
			}
		}
	}
}

func (s *Strategy) mergeMachinePSO(machines []*cloud.Machine) bool {
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
