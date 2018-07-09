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
	//fmt.Println("preDeployRestDiskNeedSkip", restDisk, maxDisk)

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
	if m.MachineId == 3163 {
		cloud.SetDebug(true)
	}

	fmt.Printf("preDeploySearch machineId=%d,disk=%d,instanceCount=%d\n",
		m.MachineId, instanceList[0].Config.Disk, len(instanceList))
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
		skipHigh := first.Config.Disk > 167 && instance.Config.Disk > 167

		if restDisk >= 40 {
			maxCpu := float64(0)
			for t, v := range resource.Cpu {
				cpu := v + instance.Config.Cpu[t]
				if cpu > maxCpu {
					maxCpu = cpu
				}
			}

			if maxCpu > 40 {
				skipHigh = true
				//fmt.Println("skipHigh ", maxCpu)
			}
		}

		finish := false
		for {
			if same {
				break
			}

			if skipHigh {
				lastAppId = instance.Config.AppId
				break
			}

			skipRestDisk := s.preDeployRestDiskNeedSkip(restDisk, instance.Config.Disk)
			if skipRestDisk {
				lastAppId = instance.Config.AppId
				break
			}

			resourceLimitOK := cloud.ConstraintCheckResourceLimit(resource, &instance.Config.Resource, limit, cloud.MaxCpuRatio)
			if !resourceLimitOK {
				lastAppId = instance.Config.AppId
				break
			}

			appInferenceOK := cloud.ConstraintCheckAppInterferenceAddInstance(
				instance.Config.AppId, appCount, s.R.AppInterferenceConfigMap)
			if !appInferenceOK {
				lastAppId = instance.Config.AppId
				break
			}

			if cloud.DebugEnabled {
				//fmt.Printf("preDeploySearch add depth=%d,i=%d,disk=%d,iDisk=%d,\n",
				//	depth, i, resource.Disk, instance.Config.Disk)
			}

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
						finish = true
					}
				}

				count++
				if count > PreDeploySearchCountMax && len(bestResult) > 1 {
					finish = true
				}
			} else {
				resource.AddResource(&instance.Config.Resource)
				appCount.Add(instance.Config.AppId)
				lastAppId = 0
				if cloud.DebugEnabled {
					//fmt.Printf("preDeploySearch depth++ depth=%d,offset=%d,disk=%d\n", depth, i, instance.Config.Disk)
				}

				instanceStack[depth] = instance
				offsetStack[depth] = i
				depth++
			}
		}

		if finish {
			break
		}

		for {
			if i < size-1 {
				i++
				break
			}

			lastAppId = 0

			//fmt.Println("depth--")

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

	//fmt.Println(offsetStack)
	//fmt.Println(depth)
	//for _, v := range instanceStack[:depth] {
	//	fmt.Println(v.Config.Disk)
	//}

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
			if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
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

//机器按照磁盘从大到小排序，再按平均内存加磁盘从小到大排序
func (s *BestFitStrategy) sortMachineByDiskDescCpuMem(machines []*cloud.Machine, cpuMemAsc bool) {
	for _, v := range machines {
		v.CalcTimedResourceStatistics()
	}
	sort.Slice(machines, func(i, j int) bool {
		if machines[i].Disk > machines[j].Disk {
			return true
		} else if machines[i].Disk == machines[j].Disk {
			a1 := machines[i].CpuAvg*(float64(288)/float64(46)) + machines[i].MemAvg
			a2 := machines[j].CpuAvg*(float64(288)/float64(46)) + machines[j].MemAvg
			if cpuMemAsc {
				if a1 < a2 {
					return true
				} else {
					return false
				}
			} else {
				if a1 > a2 {
					return true
				} else {
					return false
				}
			}

		} else {
			return false
		}
	})
}

func (s *BestFitStrategy) analysisMachineDiskDistribution(machines []*cloud.Machine) {
	diskCounts := make([]*cloud.DiskCount, 0)
	for _, v := range machines {
		has := false
		for _, p := range diskCounts {
			if p.Disk == v.Disk {
				p.Count++
				has = true
			}
		}
		if !has {
			diskCounts = append(diskCounts, &cloud.DiskCount{Disk: v.Disk, Count: 1})
		}
	}
	sort.Slice(diskCounts, func(i, j int) bool {
		return diskCounts[i].Disk > diskCounts[j].Disk
	})
	fmt.Println("fillMachines diskCounts")
	for _, v := range diskCounts {
		fmt.Println(v.Disk, v.Count)
	}
}

func (s *BestFitStrategy) sortInstanceByCpuMem(instances []*cloud.Instance, asc bool) {
	sort.Slice(instances, func(i, j int) bool {
		a1 := instances[i].Config.CpuAvg*(float64(288)/float64(46)) + instances[i].Config.MemAvg
		a2 := instances[j].Config.CpuAvg*(float64(288)/float64(46)) + instances[j].Config.MemAvg
		order := false
		if asc {
			order = a1 < a2
		} else {
			order = a1 > a2
		}

		if order {
			return true
		} else if a1 == a2 {
			if instances[i].Config.AppId < instances[j].Config.AppId {
				return true
			} else {
				return false
			}
		} else {
			return false
		}

	})
}

func (s *BestFitStrategy) removeInstances(instances []*cloud.Instance, removes []*cloud.Instance) (rest []*cloud.Instance) {
	rest = make([]*cloud.Instance, 0)
	for _, v := range instances {
		has := false
		for _, i := range removes {
			if i.InstanceId == v.InstanceId {
				has = true
				break
			}
		}
		if !has {
			rest = append(rest, v)
		}
	}

	return rest
}

func (s *BestFitStrategy) instanceContains(instances []*cloud.Instance, instanceId int) bool {
	for _, v := range instances {
		if v.InstanceId == instanceId {
			return true
		}
	}

	return false
}

func (s *BestFitStrategy) fillMachine(m *cloud.Machine, i80 []*cloud.Instance, i60 []*cloud.Instance, i40 []*cloud.Instance) (
	i80Rest []*cloud.Instance, i60Rest []*cloud.Instance, i40Rest []*cloud.Instance, err error) {

	i80Deployed := make([]*cloud.Instance, 32)
	i80DeployedCount := 0
	i60Deployed := make([]*cloud.Instance, 32)
	i60DeployedCount := 0
	i40Deployed := make([]*cloud.Instance, 32)
	i40DeployedCount := 0

	//todo 降低同一app在同一机器的部署

	disk := 1024 - m.Disk
	disk = disk - disk%20
	if disk < 40 {
		return i80, i60, i40, err
	}

	if len(i40) > 0 { //有40的先用40的
		need60 := false
		if disk%40 == 20 { //需要补一个60
			need60 = true
			disk -= 60
		}

		for _, instance := range i40 { //上面去掉了60，这里的disk是40的倍数
			if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
				continue
			}

			m.AddInstance(instance)
			i40Deployed[i40DeployedCount] = instance
			i40DeployedCount++
			disk -= 40
			fmt.Printf("fillMachine 40 disk=%d\n", disk)
			if disk == 0 {
				break
			}
		}
		if disk != 0 {
			if disk >= 80 {
				for _, instance80 := range i80 {
					if disk == 120 { //后面补2个60
						break
					}

					if !m.ConstraintCheck(instance80, cloud.MaxCpuRatio) {
						continue
					}

					fmt.Println("fillMachine 40 80", disk)
					m.AddInstance(instance80)
					i80Deployed[i80DeployedCount] = instance80
					i80DeployedCount++
					disk -= 80
					if disk == 0 {
						break
					}
				}
			}

			if disk == 40 { //此时应拿出3个40，补两个80
				fmt.Println("fillMachine 40 disk=40 此时应拿出3个40，补两个80")
				for i := 0; i < 3; i++ {
					m.RemoveInstance(i40Deployed[i40DeployedCount-1].InstanceId)
					i40DeployedCount--
				}
				for i := 0; i < 2; i++ {
					deployed := false
					for _, instance := range i80 {
						if s.instanceContains(i80Deployed[:i80DeployedCount], instance.InstanceId) {
							continue
						}

						if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
							continue
						}

						fmt.Println("fillMachine 80 60")
						m.AddInstance(instance)
						i80Deployed[i80DeployedCount] = instance
						i80DeployedCount++
						deployed = true
						break
					}
					if !deployed {
						return nil, nil, nil, fmt.Errorf("fillMachine failed 40 disk=40 %d", m.MachineId)
					}
				}
			} else if disk == 120 { //补两个60
				fmt.Println("fillMachine disk == 120 补两个60")
				for i := 0; i < 2; i++ {
					deployed := false
					for _, instance := range i60 {
						if s.instanceContains(i60Deployed[:i60DeployedCount], instance.InstanceId) {
							continue
						}

						if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
							continue
						}

						fmt.Println("fillMachine 80 60")
						m.AddInstance(instance)
						i60Deployed[i60DeployedCount] = instance
						i60DeployedCount++
						deployed = true
						break
					}
					if !deployed {
						return nil, nil, nil, fmt.Errorf("fillMachine failed disk == 120 %d", m.MachineId)
					}
				}
			}
		}

		if need60 {
			for _, instance := range i60 {
				if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
					continue
				}

				fmt.Println("fillMachine 60")
				m.AddInstance(instance)
				i60Deployed[i60DeployedCount] = instance
				i60DeployedCount++
			}
		}
	} else {
		if len(i80) > 0 { //有80的再用80的
			need60 := 0
			if disk%80 == 20 {
				need60 = 3
				if disk < 180 {
					panic(fmt.Errorf("need60=3 %d", disk))
				}
				disk -= 180
			} else if disk%80 == 40 {
				need60 = 2
				if disk < 120 {
					panic(fmt.Errorf("need60=2 %d", disk))
				}
				disk -= 120
			} else if disk%80 == 60 {
				need60 = 1
				disk -= 60
			}

			fmt.Println("80 need60", need60)

			for _, instance := range i80 {
				if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
					continue
				}

				m.AddInstance(instance)
				i80Deployed[i80DeployedCount] = instance
				i80DeployedCount++
				disk -= 80
				fmt.Printf("fillMachine 80 disk=%d\n", disk)
				if disk == 0 {
					break
				}
			}

			if disk > 0 { //这里是80的整数倍
				fmt.Println("fill 80 disk>0 这里是80的整数倍 disk=", disk)
				for _, instance := range i60 {
					if s.instanceContains(i60Deployed[:i60DeployedCount], instance.InstanceId) {
						continue
					}

					if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
						continue
					}

					m.AddInstance(instance)
					i60Deployed[i60DeployedCount] = instance
					i60DeployedCount++
					disk -= 60
					fmt.Printf("fillMachine 80 disk=%d\n", disk)
					if disk < 60 {
						break
					}
				}

				if disk > 0 { //这里可能是20或40，20不用考虑，如果是40，拿走一个80补两个60，80放到另一个里，还是会多40
					fmt.Println("fill 80 disk>0 这里是80的整数倍 restDisk=", disk)
				}
			}

			if need60 > 0 {
				for i := 0; i < need60; i++ {
					deployed := false
					for _, instance := range i60 {
						if s.instanceContains(i60Deployed[:i60DeployedCount], instance.InstanceId) {
							continue
						}

						if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
							continue
						}

						fmt.Println("fillMachine 80 60")
						m.AddInstance(instance)
						i60Deployed[i60DeployedCount] = instance
						i60DeployedCount++
						deployed = true
						break
					}
					if !deployed {
						return nil, nil, nil, fmt.Errorf("fillMachine failed 80,60 %d", m.MachineId)
					}
				}
			}
		} else { //最后用60的补满
			//cloud.SetDebug(true)
			for _, instance := range i60 {
				if s.instanceContains(i60Deployed[:i60DeployedCount], instance.InstanceId) {
					continue
				}

				if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
					continue
				}

				cpuMax := float64(0)
				for _, v := range m.Cpu {
					if v > cpuMax {
						cpuMax = v
					}
				}
				if cpuMax >= float64(32) {
					cpuRest := float64(46) - cpuMax
					count := float64(disk) / float64(60)
					cpuHigh := false
					for _, v := range instance.Config.Cpu {
						if v > cpuRest/count {
							cpuHigh = true
							break
						}
					}
					if cpuHigh {
						//fmt.Println("fill cpu high ", cpuMax)
						continue
					}
				}

				m.AddInstance(instance)

				i60Deployed[i60DeployedCount] = instance
				i60DeployedCount++
				disk -= 60
				fmt.Printf("fillMachine 60 disk=%d\n", disk)
				if disk == 0 { //肯定是60的整数倍
					break
				}
			}

			if disk != 0 {
				return nil, nil, nil, fmt.Errorf("fill 60 disk!=0 %d ", disk)
			}
		}
	}

	m.Resource.DebugPrint()

	return s.removeInstances(i80, i80Deployed[:i80DeployedCount]),
		s.removeInstances(i60, i60Deployed[:i60DeployedCount]),
		s.removeInstances(i40, i40Deployed[:i40DeployedCount]), err
}

func (s *BestFitStrategy) fillMachines(machines []*cloud.Machine, instances []*cloud.Instance) (restInstances []*cloud.Instance, err error) {
	s.analysisMachineDiskDistribution(machines)
	//cloud.AnalysisDiskDistributionByInstance(instances)

	//机器按照磁盘从大到小，cpu内存从小到大排序
	s.sortMachineByDiskDescCpuMem(machines, true)

	//实例按照磁盘大小分组
	i80Rest := make([]*cloud.Instance, 0)
	i60Rest := make([]*cloud.Instance, 0)
	i40Rest := make([]*cloud.Instance, 0)
	for _, v := range instances {
		if v.Config.Disk == 80 {
			i80Rest = append(i80Rest, v)
		} else if v.Config.Disk == 60 {
			i60Rest = append(i60Rest, v)
		} else if v.Config.Disk == 40 {
			i40Rest = append(i40Rest, v)
		}
	}
	s.sortInstanceByCpuMem(i80Rest, false)
	s.sortInstanceByCpuMem(i60Rest, false)
	s.sortInstanceByCpuMem(i40Rest, false)

	for i, v := range i80Rest {
		v.Config.DebugPrint()
		if i > 10 {
			break
		}
	}

	for i, v := range i60Rest {
		v.Config.DebugPrint()
		if i > 10 {
			break
		}
	}

	for i, v := range i40Rest {
		v.Config.DebugPrint()
		if i > 10 {
			break
		}
	}

	//使用磁盘40补齐所有机器到60的倍数
	i40Deployed := make(map[int]*cloud.Instance)
	for _, m := range machines {
		disk := 1024 - m.Disk
		disk = disk - disk%20
		if disk < 40 {
			continue
		}

		need40 := 0
		if disk%60 == 0 {
			need40 = 0
		} else if disk%60 == 20 { //补两个40
			need40 = 2
		} else if disk%60 == 40 { //补1个40
			need40 = 1
		}

		if need40 == 0 {
			continue
		}

		for i := 0; i < need40; i++ {
			deployed := false
			for _, instance := range i40Rest {
				_, has := i40Deployed[instance.InstanceId]
				if has {
					continue
				}

				if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
					continue
				}

				m.AddInstance(instance)
				i40Deployed[instance.InstanceId] = instance
				deployed = true
				break
			}
			if !deployed {
				return nil, fmt.Errorf("fillMachines %%60 !deployed machineId=%d", m.MachineId)
			}
		}
	}
	i40Rest = s.preDeployRemoveDeployed(i40Rest, i40Deployed)
	fmt.Printf("fillMachines rest %d,deployed=%d\n", len(i40Rest), len(i40Deployed))

	//补满所有已部署的机器
	n := 0
	for i, m := range machines {
		if i >= 424 {
			//cloud.SetDebug(true)
		}

		fmt.Println("fillMachines m.Disk", m.Disk)
		i80Rest, i60Rest, i40Rest, err = s.fillMachine(m, i80Rest, i60Rest, i40Rest)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		fmt.Println("fillMachines", len(i80Rest), len(i60Rest), len(i40Rest))

		if len(i40Rest) == 1 {
			n++
			if n == 20 {
				break
			}
		}

		if i >= 1000 {
			//break
		}
	}

	restInstances = make([]*cloud.Instance, 0)
	restInstances = append(restInstances, i80Rest...)
	restInstances = append(restInstances, i60Rest...)
	restInstances = append(restInstances, i40Rest...)

	return restInstances, nil
}

func (s *BestFitStrategy) preDeployBigDisk(instanceList []*cloud.Instance, machineList []*cloud.Machine) (
	restInstances []*cloud.Instance, restMachines []*cloud.Machine, err error) {

	// 部署磁盘大于100的，每个实例一个机器
	// 磁盘167的，每个磁盘两个
	// 磁盘150的，用于填充之前的单数［650，250，167*2=334］机器，剩下的每两个一个机器
	// 磁盘100的，每个机器5个，尾数再分配一个机器
	deployedInstances := make(map[int]*cloud.Instance)
	deployedMachines := make([]*cloud.Machine, 0)
	machineIndex := 0
	instances150 := make([]*cloud.Instance, 0)
	for _, instance := range instanceList {
		if instance.Config.Disk < 100 {
			continue
		}

		if instance.Config.Disk == 150 { //150的之后再补
			instances150 = append(instances150, instance)
			continue
		}

		if instance.Config.Disk == 167 {
			deployed167 := false
			for _, m := range deployedMachines {
				if m.InstanceArrayCount == 1 && m.InstanceArray[0].Config.Disk == 167 {
					if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
						return nil, nil, fmt.Errorf("preDeployBigDisk.167 ConstraintCheck failed")
					}
					s.R.CommandDeployInstance(instance, m)
					deployedInstances[instance.InstanceId] = instance
					deployed167 = true
					break
				}
			}
			if !deployed167 {
				s.R.CommandDeployInstance(instance, machineList[machineIndex])
				deployedInstances[instance.InstanceId] = instance
				deployedMachines = append(deployedMachines, machineList[machineIndex])
				machineIndex++
			}
		} else if instance.Config.Disk == 100 {
			deployed100 := false
			for _, m := range deployedMachines {
				if m.InstanceArrayCount >= 5 {
					continue
				}
				is100 := false
				for _, v := range m.InstanceArray[:m.InstanceArrayCount] {
					if v.Config.Disk == 100 {
						is100 = true
						break
					}
				}
				if !is100 {
					continue
				}

				if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
					continue
				}
				s.R.CommandDeployInstance(instance, m)
				deployedInstances[instance.InstanceId] = instance
				deployed100 = true
				break
			}
			if !deployed100 {
				s.R.CommandDeployInstance(instance, machineList[machineIndex])
				deployedInstances[instance.InstanceId] = instance
				deployedMachines = append(deployedMachines, machineList[machineIndex])
				machineIndex++
			}
		} else {
			s.R.CommandDeployInstance(instance, machineList[machineIndex])
			deployedInstances[instance.InstanceId] = instance
			deployedMachines = append(deployedMachines, machineList[machineIndex])
			machineIndex++
		}
	}

	//磁盘150的按照CPU内存从大到小排序
	sort.Slice(instances150, func(i, j int) bool {
		a1 := instances150[i].Config.CpuAvg*(float64(288)/float64(46)) + instances150[i].Config.MemAvg
		a2 := instances150[j].Config.CpuAvg*(float64(288)/float64(46)) + instances150[j].Config.MemAvg
		if a1 > a2 {
			return true
		} else {
			return false
		}
	})

	//磁盘150插入已部署的单数机器
	s.sortMachineByDiskDescCpuMem(deployedMachines, true)
	for _, m := range deployedMachines {
		if m.Disk == 650 || m.Disk == 334 || m.Disk == 250 {
			deployed150 := false
			for _, instance := range instances150 {
				_, has := deployedInstances[instance.InstanceId]
				if has {
					continue
				}

				if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
					continue
				}

				s.R.CommandDeployInstance(instance, m)
				deployedInstances[instance.InstanceId] = instance
				deployed150 = true
				break
			}
			if !deployed150 {
				return nil, nil, fmt.Errorf("deployed150 failed,machineId=%d,disk=%d",
					m.MachineId, m.Disk)
			}
		}
	}

	//磁盘150剩下的两个一组插入
	for _, instance := range instances150 {
		_, has := deployedInstances[instance.InstanceId]
		if has {
			continue
		}

		deployed150 := false
		for _, m := range deployedMachines {
			if m.InstanceArrayCount == 1 && m.InstanceArray[0].Config.Disk == 150 {
				if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
					return nil, nil, fmt.Errorf("preDeployBigDisk.150 ConstraintCheck failed")
				}
				s.R.CommandDeployInstance(instance, m)
				deployedInstances[instance.InstanceId] = instance
				deployed150 = true
				break
			}
		}
		if !deployed150 {
			s.R.CommandDeployInstance(instance, machineList[machineIndex])
			deployedInstances[instance.InstanceId] = instance
			deployedMachines = append(deployedMachines, machineList[machineIndex])
			machineIndex++
		}
	}

	//填充这些机器
	restInstances = s.preDeployRemoveDeployed(instanceList, deployedInstances)
	restInstances, err = s.fillMachines(deployedMachines, restInstances)
	if err != nil {
		return nil, nil, err
	}

	restMachines = make([]*cloud.Machine, 0)
	for _, m := range machineList {
		has := false
		for _, v := range deployedMachines {
			if v.MachineId == m.MachineId {
				has = true
				break
			}
		}
		if !has {
			restMachines = append(restMachines, m)
		}
	}

	return restInstances, restMachines, nil
}
