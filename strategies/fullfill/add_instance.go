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
	cpu := m.CpuMax / (m.LevelConfig.Cpu * cloud.MaxCpuRatio)
	mem := m.MemMax / m.LevelConfig.Mem

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

	return typ, d
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

		offset = int(float64(len(pool)) * 2 * d)
		if len(pool)-offset < 10000 {
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

			derivation := m.CalcDeviationWithInstance(instance)
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
				break
			}
		}
	}

	return cloud.InstancesRemove(instances, deployed[:deployedCount]),
		cloud.InstancesRemove(instancesByDisk, deployed[:deployedCount]),
		cloud.InstancesRemove(instancesByCpu, deployed[:deployedCount]),
		cloud.InstancesRemove(instancesByMem, deployed[:deployedCount]),
		nil
}
