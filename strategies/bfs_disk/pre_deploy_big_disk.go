package bfs_disk

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

const minDisk = 40

func (s *BestFitStrategy) preDeployDistance(
	r *cloud.Resource, instance *cloud.Instance, limit *cloud.MachineLevelConfig, log bool) float64 {

	cpu := float64(0)
	for i, v := range r.Cpu {
		cpu += v + instance.Config.Cpu[i]
	}
	cpu = cpu / float64(len(r.Cpu))
	rCpu := (limit.Cpu - cpu) / limit.Cpu

	mem := float64(0)
	for i, v := range r.Mem {
		mem += v + instance.Config.Mem[i]
	}
	mem = mem / float64(len(r.Mem))
	rMem := (limit.Mem - mem) / limit.Mem

	if log {
		fmt.Printf("preDeployDistance %f %f %f %f\n", cpu, rCpu, mem, rMem)
	}

	return rCpu + rMem
}

func (s *BestFitStrategy) preDeployRestDiskNeedSkip(restDisk int, maxDisk int) bool {
	if restDisk < 10 {
		return false
	}

	if restDisk < 40 {
		return true
	}

	switch maxDisk {
	case 40:
		{
			if restDisk%40 < 10 {
				return false
			} else {
				return true
			}
		}
	case 60:
		{
			if restDisk%20 < 10 {
				return false
			} else {
				return true
			}
		}
	case 80:
		{
			if restDisk%20 < 10 {
				return false
			} else {
				return true
			}
		}
	case 100:
		{
			if restDisk%20 < 10 {
				return false
			} else {
				return true
			}
		}
	case 120:
		{
			if restDisk%20 < 10 {
				return false
			} else {
				return true
			}
		}
	default:
		return false
	}

	return false
}

func (s *BestFitStrategy) preDeploySearch(m *cloud.Machine, instanceList []*cloud.Instance) ([]*cloud.Instance, error) {
	if m.MachineId == 3174 {
		cloud.SetDebug(true)
	}

	fmt.Println(len(instanceList))

	fmt.Printf("preDeploySearch machineId=%d,disk=%d\n", m.MachineId, instanceList[0].Config.Disk)
	limit := &cloud.MachineLevelConfig{}
	limit.Cpu = m.LevelConfig.Cpu / float64(2)
	limit.Mem = m.LevelConfig.Mem
	limit.Disk = m.LevelConfig.Disk
	limit.P = m.LevelConfig.P
	limit.M = m.LevelConfig.M
	limit.PM = m.LevelConfig.PM

	first := instanceList[0]
	if limit.Disk-first.Config.Disk < minDisk {
		//fmt.Printf("preDeploySearch limit.Disk-first.Config.Disk < minDisk  %d %d\n", limit.Disk, first.Config.Disk)
		return []*cloud.Instance{first}, nil
	}

	resource := &cloud.Resource{}
	resource.AddResource(&first.Config.Resource)
	appCount := cloud.NewAppCountCollection()
	appCount.Add(first.Config.AppId)

	var bestResult []*cloud.Instance
	bestDistance := float64(2)

	depth := 0
	instanceStack := make([]*cloud.Instance, 32)
	instanceStack[depth] = first
	offsetStack := make([]int, 32)
	offsetStack[depth] = 0
	depth++
	offsetStack[depth] = 1

	size := len(instanceList)
	lastAppId := 0
	count := 0
	for i := offsetStack[depth]; i < size; {
		instance := instanceList[i]
		same := lastAppId != 0 && instance.Config.AppId == lastAppId
		restDisk := limit.Disk - (resource.Disk + instance.Config.Disk)
		if !same && !s.preDeployRestDiskNeedSkip(restDisk, instance.Config.Disk) &&
			cloud.ConstraintCheckResourceLimit(resource, &instance.Config.Resource, limit) &&
			cloud.ConstraintCheckAppInterferenceAddInstance(instance.Config.AppId, appCount, s.R.AppInterferenceConfigMap) {

			if cloud.DebugEnabled {
				fmt.Printf("preDeploySearch add depth=%d,i=%d,disk=%d,iDisk=%d,\n",
					depth, i, resource.Disk, instance.Config.Disk)
			}

			lastAppId = instance.Config.AppId
			if restDisk < 10 {
				distance := s.preDeployDistance(resource, instance, limit, false)

				if cloud.DebugEnabled {
					fmt.Printf("preDeploySearch %d,%d %f %f\n", resource.Disk, instance.Config.Disk, distance, bestDistance)
				}

				if distance < bestDistance-0.000001 {
					fmt.Printf("preDeploySearch best %d,%d %f\n",
						resource.Disk, instance.Config.Disk, distance)
					//resource.DebugPrint()
					bestDistance = distance
					bestResult = make([]*cloud.Instance, 0)
					bestResult = append(bestResult, instanceStack[:depth]...)
					bestResult = append(bestResult, instance)

					if distance < 0.1 {
						break
					}
				}

				count++
				if count > PreDeploySearchCountMax && len(bestResult) > 1 {
					break
				}
			} else {
				resource.AddResource(&instance.Config.Resource)
				appCount.Add(instance.Config.AppId)
				lastAppId = 0
				if cloud.DebugEnabled {
					fmt.Printf("preDeploySearch depth++ depth=%d,offset=%d,disk=%d\n", depth, i, instance.Config.Disk)
				}

				instanceStack[depth] = instance
				offsetStack[depth] = i
				depth++
			}
		}

		for {
			if i < size-1 {
				i++
				break
			}

			lastAppId = 0

			depth--
			if depth == 0 {
				break
			}

			ins := instanceStack[depth]
			if ins != nil {
				resource.RemoveResource(&ins.Config.Resource)
				appCount.Remove(ins.Config.AppId)
			}

			i = offsetStack[depth] + 1
			if i == size-1 || i == size {
				fmt.Printf("preDeploySearch i = offsetStack[depth] + 1 i == size-1\n")
				fmt.Println(resource.Disk, depth, offsetStack)
				for _, v := range instanceStack[:depth+1] {
					fmt.Println(v.Config.Disk)
				}
				cloud.AnalysisDiskDistributionByInstance(instanceList)
			} else {
				break
			}
		}

		if depth == 0 {
			fmt.Printf("preDeploySearch finished\n")
			break
		}
	}

	fmt.Println(offsetStack)
	fmt.Println(depth)
	for _, v := range instanceStack[:depth] {
		fmt.Println(v.Config.Disk)
	}

	fmt.Printf("preDeploySearch bestResult=%f,count=%d\n", bestDistance, len(bestResult))

	return bestResult, nil
}

func (s *BestFitStrategy) preDeployFill(m *cloud.Machine, instanceList []*cloud.Instance) ([]*cloud.Instance, error) {
	deployed := make(map[int]*cloud.Instance)

	result, err := s.preDeploySearch(m, instanceList)
	if err != nil {
		return nil, err
	}

	if result != nil {
		for _, instance := range result {
			if !m.ConstraintCheck(instance) {
				return nil, fmt.Errorf("BestFitStrategy.preDeployFill ConstraintCheck failed,machineId=%d,instanceId=%d",
					m.MachineId, instance.InstanceId)
			}

			s.R.CommandDeployInstance(instance, m)
			deployed[instance.InstanceId] = instance
		}

		//m.DebugPrint()
	}

	return s.preDeployRemoveDeployed(instanceList, deployed), nil
}

func (s *BestFitStrategy) preDeployBigDisk(instanceList []*cloud.Instance) (
	restInstances []*cloud.Instance, restMachines []*cloud.Machine, err error) {
	sort.Slice(instanceList, func(i, j int) bool {
		if instanceList[i].Config.Disk > instanceList[j].Config.Disk {
			return true
		} else if instanceList[i].Config.Disk == instanceList[j].Config.Disk {
			a1 := instanceList[i].Config.CpuAvg*(float64(46)/float64(288)) + instanceList[i].Config.MemAvg
			a2 := instanceList[j].Config.CpuAvg*(float64(46)/float64(288)) + instanceList[j].Config.MemAvg
			if a1 > a2 {
				return true
			} else if a1 == a2 {
				if instanceList[i].Config.AppId < instanceList[j].Config.AppId {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	})

	restInstances = make([]*cloud.Instance, 0)
	restInstances = append(restInstances, instanceList...)

	for i, m := range s.machineDeployList {
		if m.InstanceArrayCount != 0 {
			continue
		}

		if i >= PreDeployMachineCount {
			break
		}

		restInstances, err = s.preDeployFill(m, restInstances)
		if err != nil {
			return nil, nil, err
		}
	}

	restMachines = make([]*cloud.Machine, 0)
	for _, v := range s.machineDeployList {
		if v.InstanceArrayCount == 0 {
			restMachines = append(restMachines, v)
		}
	}

	fmt.Printf("BestFitStrategy.preDeployBigDisk instantCount=%d,machineCount=%d\n",
		len(instanceList)-len(restInstances), len(s.machineDeployList)-len(restMachines))

	return restInstances, restMachines, nil
}
