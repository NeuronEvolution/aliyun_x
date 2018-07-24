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

	for i := 0; ; i++ {
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
			return nil
		}

		if m.LevelConfig.Disk == 1024 {
			if m.Disk <= 980 {
				fmt.Println(i)
				m.DebugPrint()
			}
		} else {
			if m.Disk <= 560 {
				fmt.Println(i)
				m.DebugPrint()
			}
		}

		if i == 3000 {
			//for _, v := range instancesByCpu {
			//v.Config.DebugPrint()
			//}
		}

		fmt.Printf("AddInstanceList restInstances count %d,%d,%d\n", i, len(restInstances), m.Disk)
		//m.DebugPrint()
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

	return typ, math.Pow(max+0.6, 10) / 50
}

func (s *Strategy) measureWithInstance(m *cloud.Machine, instance *cloud.Instance) (d float64) {
	cpuMax := float64(0)
	for i, v := range m.Cpu {
		cpu := v + instance.Config.Cpu[i]
		if cpu > cpuMax {
			cpuMax = cpu
		}
	}

	memMax := float64(0)
	for i, v := range m.Mem {
		mem := v + instance.Config.Mem[i]
		if mem > memMax {
			memMax = mem
		}
	}

	cpu := cpuMax / (m.LevelConfig.Cpu * cloud.MaxCpuRatio)
	mem := memMax / m.LevelConfig.Mem
	disk := float64(m.Disk+instance.Config.Disk) / float64(m.LevelConfig.Disk)

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

	if !m.ConstraintCheck(instances[0], cloud.MaxCpuRatio) {
		return nil, nil, nil, nil,
			fmt.Errorf("fillMachine ConstraintCheck instances[0] failed")
	}

	m.AddInstance(instances[0])
	deployed[deployedCount] = instances[0]
	deployedCount++

	trySameApp := false
	resetOffset := false

	for {
		has := false

		var pool []*cloud.Instance
		offset := 0
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

		offset = int(float64(len(pool)) * d)
		n := len(instances)
		if n > 60000 {
			if offset > len(pool)*95/100 {
				offset = len(pool) * 95 / 100
			}
		} else if n > 50000 {
			if offset > len(pool)*95/100 {
				offset = len(pool) * 95 / 100
			}
		} else if n > 40000 {
			if offset > len(pool)*95/100 {
				offset = len(pool) * 95 / 100
			}
		} else if n > 30000 {
			if offset > len(pool)*95/100 {
				offset = len(pool) * 95 / 100
			}
		} else if n > 20000 {
			if offset > len(pool)*95/100 {
				offset = len(pool) * 95 / 100
			}
		} else if n > 10000 {
			if offset > len(pool)*95/100 {
				offset = len(pool) * 95 / 100
			}
		} else {
			offset = 0
		}

		if resetOffset {
			offset = 0
		}

		fmt.Printf("%d %d %f\n", typ, offset, d)

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
			if derivation < minDerivation {
				if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
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
			if !trySameApp {
				trySameApp = true
			} else {
				if !resetOffset {
					resetOffset = true
				} else {
					break
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
