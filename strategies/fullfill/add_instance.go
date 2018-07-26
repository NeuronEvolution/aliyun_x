package fullfill

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
)

const TypeDisk = 1
const TypeCpu = 2
const TypeMem = 3

func (s *Strategy) AddInstanceList(instances []*cloud.Instance) (err error) {
	restInstances := cloud.InstancesCopy(instances)
	cloud.SortInstanceByTotalMax(restInstances)
	instancesByDisk := cloud.InstancesCopy(instances)
	cloud.SortInstanceByDisk(instancesByDisk)
	instancesByCpu := cloud.InstancesCopy(instances)
	cloud.SortInstanceByCpu(instancesByCpu)
	instancesByMem := cloud.InstancesCopy(instances)
	cloud.SortInstanceByMem(instancesByMem)

	n := 0
	for i := 0; i < 15500; i++ {
		m := s.R.MachineFreePool.PeekMachine()
		if m == nil {
			return fmt.Errorf("AddInstanceList PeekMachine no machine")
		}

		restInstances, instancesByDisk, instancesByCpu, instancesByMem, err =
			s.fillMachine(m, restInstances, instancesByDisk, instancesByCpu, instancesByMem)
		if err != nil {
			return err
		}

		if len(restInstances) == 0 {
			break
		}

		if m.LevelConfig.Disk == 1024 {
			if m.Disk <= 980 {
				//fmt.Println(i)
				//m.DebugPrint()
				n++
			}
		} else {
			if m.Disk < 560 {
				//fmt.Println(i)
				//m.DebugPrint()
				n++
			}
		}

		if i == 3000 {
			//for _, v := range instancesByCpu {
			//v.Config.DebugPrint()
			//}
		}

		//fmt.Printf("AddInstanceList restInstances count %d,%d,%d\n", i, len(restInstances), m.Disk)

		m.Resource.DebugPrint()
	}

	fmt.Println("AddInstanceList rest", len(restInstances), n)

	deployed := make([]*cloud.Instance, len(restInstances))
	deployedCount := 0
	for _, instance := range restInstances {
		instance.Config.DebugPrint()
		err = s.forceAddInstance(instance)
		if err != nil {
			fmt.Println("forceAddInstance failed", err)
			continue
		}

		deployed[deployedCount] = instance
		deployedCount++
	}

	fmt.Println(len(restInstances), len(instancesByDisk), len(instancesByCpu), len(instancesByMem))

	restInstances = cloud.InstancesRemove(restInstances, deployed[:deployedCount])
	instancesByDisk = cloud.InstancesRemove(instancesByDisk, deployed[:deployedCount])
	instancesByCpu = cloud.InstancesRemove(instancesByCpu, deployed[:deployedCount])
	instancesByMem = cloud.InstancesRemove(instancesByMem, deployed[:deployedCount])

	fmt.Println(len(restInstances), len(instancesByDisk), len(instancesByCpu), len(instancesByMem))

	for i := 0; ; i++ {
		if len(restInstances) == 0 {
			break
		}

		m := s.R.MachineFreePool.PeekMachine()
		if m == nil {
			return fmt.Errorf("AddInstanceList PeekMachine no machine")
		}

		restInstances, instancesByDisk, instancesByCpu, instancesByMem, err =
			s.fillMachine(m, restInstances, instancesByDisk, instancesByCpu, instancesByMem)
		if err != nil {
			return err
		}

		m.Resource.DebugPrint()
		//fmt.Printf("AddInstanceList again restInstances count %d,%d,%d\n", i, len(restInstances), m.Disk)
	}

	return nil
}

func (s *Strategy) measureTooHigh(m *cloud.Machine) (typ int, d float64) {
	disk := float64(m.Disk) / float64(m.LevelConfig.Disk)
	cpuMax := float64(0)
	for _, v := range m.Cpu {
		if v > cpuMax {
			cpuMax = v
		}
	}
	memMax := float64(0)
	for _, v := range m.Mem {
		if v > memMax {
			memMax = v
		}
	}

	cpu := cpuMax / (m.LevelConfig.Cpu * cloud.MaxCpuRatio)
	mem := memMax / m.LevelConfig.Mem

	max := float64(0)
	if disk > max {
		typ = TypeDisk
		max = disk
	}

	if cpu > max {
		typ = TypeCpu
		max = cpu
	}

	if mem > max {
		typ = TypeMem
		max = mem
	}

	switch typ {
	case TypeDisk:
		d = (disk - cpu) + (disk - mem)
	case TypeCpu:
		d = (cpu - disk) + (cpu - mem)
	case TypeMem:
		d = (mem - disk) + (mem - cpu)
	default:
		break
	}

	return typ, math.Pow(max+1, 10) / 780
}

func (s *Strategy) measureWithInstance(m *cloud.Machine, instance *cloud.Instance) (d float64) {
	disk := float64(m.Disk+instance.Config.Disk) / float64(m.LevelConfig.Disk)
	if disk > 1 {
		return disk
	}

	cpuMax := float64(0)
	for i, v := range m.Cpu {
		cpu := v + instance.Config.Cpu[i]
		if cpu > cpuMax {
			cpuMax = cpu
		}
	}
	cpu := cpuMax / (m.LevelConfig.Cpu * cloud.MaxCpuRatio)
	if cpu > 1 {
		return cpuMax
	}

	memMax := float64(0)
	for i, v := range m.Mem {
		mem := v + instance.Config.Mem[i]
		if mem > memMax {
			memMax = mem
		}
	}
	mem := memMax / m.LevelConfig.Mem
	if mem > 1 {
		return mem
	}

	typ := 0
	max := float64(0)
	if disk > max {
		typ = TypeDisk
		max = disk
	}

	if cpu > max {
		typ = TypeCpu
		max = cpu
	}

	if mem > max {
		typ = TypeMem
		max = mem
	}

	switch typ {
	case TypeDisk:
		d = ((disk - cpu) + (disk - mem)) / 2
	case TypeCpu:
		d = ((cpu - disk) + (cpu - mem)) / 2
	case TypeMem:
		d = ((mem - disk) + (mem - cpu)) / 2
	default:
		break
	}

	return d * max
}

func (s *Strategy) fillMachine(
	m *cloud.Machine,
	instances []*cloud.Instance,
	instancesByDisk []*cloud.Instance,
	instancesByCpu []*cloud.Instance,
	instancesByMem []*cloud.Instance) (
	restInstances []*cloud.Instance,
	restInstancesByDisk []*cloud.Instance,
	restInstancesByCpu []*cloud.Instance,
	restInstancesByMem []*cloud.Instance,
	err error) {

	deployed := make([]*cloud.Instance, 32)
	deployedCount := 0

	if !m.ConstraintCheck(instances[0], 1) {
		m.Resource.DebugPrint()
		instances[0].Config.DebugPrint()
		return nil, nil, nil, nil,
			fmt.Errorf("fillMachine ConstraintCheck instances[0] failed")
	}

	m.AddInstance(instances[0])
	deployed[deployedCount] = instances[0]
	deployedCount++

	trySameApp := false
	resetOffset := false

	resetCount := 0
	offset := 0
	for {
		has := false

		var pool []*cloud.Instance

		typ, d := s.measureTooHigh(m)
		switch typ {
		case TypeDisk:
			pool = instancesByDisk
		case TypeCpu:
			pool = instancesByCpu
		case TypeMem:
			pool = instancesByMem
		default:
			pool = instances
		}

		if resetOffset {
			resetCount++
			offset = len(pool) - int(math.Pow(float64(2), float64(resetCount)))
			if offset <= 0 {
				offset = 0
			}
		} else {
			offset = int(float64(len(pool)) * d)
			if offset >= len(pool) {
				offset = len(pool) - 1
			}
		}

		if trySameApp {
			offset = 0
		}

		//fmt.Printf("%d %d %f\n", typ, offset, d)

		minDerivation := math.MaxFloat64
		for _, instance := range instances[offset:] {
			if cloud.InstancesContains(deployed[:deployedCount], instance.InstanceId) {
				continue
			}

			if !trySameApp {
				if cloud.InstancesContainsApp(deployed[:deployedCount], instance.Config.AppId) {
					continue
				}
			}

			derivation := s.measureWithInstance(m, instance)
			if derivation > 1 {
				continue
			}

			if derivation < minDerivation {
				if !m.ConstraintCheck(instance, 1) {
					continue
				}

				m.AddInstance(instance)
				deployed[deployedCount] = instance
				deployedCount++

				minDerivation = derivation
				//fmt.Println(minDerivation)
				has = true
				break
			}
		}

		if m.LevelConfig.Disk-m.Disk < 40 {
			break
		}

		if !has {
			if !resetOffset {
				resetOffset = true
			} else {
				if offset <= 0 {
					if !trySameApp {
						trySameApp = true
					} else {
						break
					}
				}
			}
		}
	}

	return cloud.InstancesRemove(instances, deployed[:deployedCount]),
		cloud.InstancesRemove(instancesByDisk, deployed[:deployedCount]),
		cloud.InstancesRemove(instancesByCpu, deployed[:deployedCount]),
		cloud.InstancesRemove(instancesByMem, deployed[:deployedCount]),
		nil
}

func (s *Strategy) forceAddInstance(instance *cloud.Instance) (err error) {
	return fmt.Errorf("forceAddInstance")
}
