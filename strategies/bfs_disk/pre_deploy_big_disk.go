package bfs_disk

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
)

//机器按照磁盘从大到小排序，再按平均内存加磁盘从小到大排序(高配机器在前)
func (s *BestFitStrategy) sortMachineByDiskDescCpuMem(machines []*cloud.Machine, cpuMemAsc bool) {
	for _, v := range machines {
		v.CalcTimedResourceStatistics()
	}
	sort.Slice(machines, func(i, j int) bool {
		if machines[i].LevelConfig.Disk > machines[j].LevelConfig.Disk {
			return true
		} else if machines[i].LevelConfig.Disk == machines[j].LevelConfig.Disk {
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

func (s *BestFitStrategy) isCpuMemHigh(m *cloud.Machine, instance *cloud.Instance, diskRest int, disk int) bool {
	cpuMax := float64(0)
	for _, v := range m.Cpu {
		if v > cpuMax {
			cpuMax = v
		}
	}

	if m.LevelConfig.Cpu == 92 {
		if cpuMax >= float64(32) {
			cpuRest := float64(46) - cpuMax
			count := float64(diskRest) / float64(disk)
			cpuHigh := false
			for _, v := range instance.Config.Cpu {
				if count > 5 {
					if v > (cpuRest * 2 / count) {
						cpuHigh = true
						break
					}
				} else {
					if v > (cpuRest * 1.2 / count) {
						cpuHigh = true
						break
					}
				}
			}
			if cpuHigh {
				if cloud.DebugEnabled {
					fmt.Println("fill cpu high ", cpuMax)
					m.Resource.DebugPrint()
					instance.Config.DebugPrint()
				}
				return true
			}
		}
	} else if m.LevelConfig.Cpu == 32 {
		if cpuMax >= float64(4) {
			cpuRest := float64(16) - cpuMax
			count := float64(diskRest) / float64(disk)
			cpuHigh := false
			if cloud.DebugEnabled {
				fmt.Println("fill", cpuRest, count, diskRest)
			}
			for _, v := range instance.Config.Cpu {
				if count > 5 {
					if v > (cpuRest * 2 / count) {
						cpuHigh = true
						break
					}
				} else {
					if v > (cpuRest * 1.2 / count) {
						cpuHigh = true
						break
					}
				}
			}
			if cpuHigh {
				if cloud.DebugEnabled {
					fmt.Println("fill cpu high ", cpuMax)
					m.Resource.DebugPrint()
					instance.Config.DebugPrint()
				}
				return true
			}
		}
	}

	memMax := float64(0)
	for _, v := range m.Mem {
		if v > memMax {
			memMax = v
		}
	}

	if m.LevelConfig.Mem == 288 {
		if memMax >= float64(144) {
			memRest := float64(288) - memMax
			count := float64(diskRest) / float64(disk)
			memHigh := false
			for _, v := range instance.Config.Mem {
				if v > memRest/count {
					memHigh = true
					break
				}
			}
			if memHigh {
				if cloud.DebugEnabled {
					fmt.Println("fill mem high ", memMax)
					m.Resource.DebugPrint()
					instance.Config.DebugPrint()
				}
				return true
			}
		}
	} else if m.LevelConfig.Mem == 64 {
		if memMax >= float64(32) {
			memRest := float64(64) - memMax
			count := float64(diskRest) / float64(disk)
			memHigh := false
			for _, v := range instance.Config.Mem {
				if v > memRest/count {
					memHigh = true
					break
				}
			}
			if memHigh {
				if cloud.DebugEnabled {
					fmt.Println("fill mem high ", memMax)
					m.Resource.DebugPrint()
					instance.Config.DebugPrint()
				}
				return true
			}
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

	disk := m.LevelConfig.Disk - m.Disk
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
			if s.isCpuMemHigh(m, instance, disk, 40) {
				continue
			}

			if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
				continue
			}

			m.AddInstance(instance)
			i40Deployed[i40DeployedCount] = instance
			i40DeployedCount++
			disk -= 40
			//fmt.Printf("fillMachine 40 disk=%d\n", disk)
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

					if s.isCpuMemHigh(m, instance80, disk, 80) {
						continue
					}

					if !m.ConstraintCheck(instance80, cloud.MaxCpuRatio) {
						continue
					}

					//fmt.Println("fillMachine 40 80", disk)
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

						if s.isCpuMemHigh(m, instance, disk, 80) {
							continue
						}

						if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
							continue
						}

						//fmt.Println("fillMachine 80 60")
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

						if s.isCpuMemHigh(m, instance, disk, 60) {
							continue
						}

						if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
							continue
						}

						//fmt.Println("fillMachine 80 60")
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
			deployed := false
			for _, instance := range i60 {
				if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
					continue
				}

				//fmt.Println("fillMachine 60")
				m.AddInstance(instance)
				i60Deployed[i60DeployedCount] = instance
				i60DeployedCount++
				deployed = true
			}
			if !deployed {
				return nil, nil, nil, fmt.Errorf("40 need60 not deployed")
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

			//fmt.Println("80 need60", need60)

			for _, instance := range i80 {
				if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
					continue
				}

				if s.isCpuMemHigh(m, instance, disk, 80) {
					continue
				}

				m.AddInstance(instance)
				i80Deployed[i80DeployedCount] = instance
				i80DeployedCount++
				disk -= 80
				//fmt.Printf("fillMachine 80 disk=%d\n", disk)
				if disk == 0 {
					break
				}
			}

			if disk > 0 { //这里是80的整数倍
				//fmt.Println("fill 80 disk>0 这里是80的整数倍 disk=", disk)
				for _, instance := range i60 {
					if s.instanceContains(i60Deployed[:i60DeployedCount], instance.InstanceId) {
						continue
					}

					if s.isCpuMemHigh(m, instance, disk, 60) {
						continue
					}

					if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
						continue
					}

					m.AddInstance(instance)
					i60Deployed[i60DeployedCount] = instance
					i60DeployedCount++
					disk -= 60
					//fmt.Printf("fillMachine 80 disk=%d\n", disk)
					if disk < 60 {
						break
					}
				}

				if disk > 0 { //这里可能是20或40，20不用考虑，如果是40，拿走一个80补两个60，80放到另一个里，还是会多40
					//fmt.Println("fill 80 disk>0 这里是80的整数倍 restDisk=", disk)
				}
			}

			if need60 > 0 {
				for i := 0; i < need60; i++ {
					deployed := false
					for _, instance := range i60 {
						if s.instanceContains(i60Deployed[:i60DeployedCount], instance.InstanceId) {
							continue
						}

						if s.isCpuMemHigh(m, instance, disk, 60) {
							continue
						}

						if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
							continue
						}

						//fmt.Println("fillMachine 80 60")
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
			lastAppId := 0
			for _, instance := range i60 {
				if instance.InstanceId == lastAppId {
					if cloud.DebugEnabled {
						fmt.Println("fill 60 instance.InstanceId == lastAppId ")
					}
					continue
				}

				if s.instanceContains(i60Deployed[:i60DeployedCount], instance.InstanceId) {
					continue
				}

				if s.isCpuMemHigh(m, instance, disk, 60) {
					continue
				}

				if !m.ConstraintCheck(instance, cloud.MaxCpuRatio) {
					continue
				}

				m.AddInstance(instance)
				lastAppId = instance.InstanceId
				i60Deployed[i60DeployedCount] = instance
				i60DeployedCount++
				disk -= 60
				if cloud.DebugEnabled {
					fmt.Printf("fillMachine 60 disk=%d\n", disk)
				}
				if disk == 0 { //肯定是60的整数倍
					break
				}
			}

			if disk != 0 {
				return nil, nil, nil, fmt.Errorf("fill 60 disk!=0 %d ", disk)
			}
		}
	}

	//m.Resource.DebugPrint()

	return s.removeInstances(i80, i80Deployed[:i80DeployedCount]),
		s.removeInstances(i60, i60Deployed[:i60DeployedCount]),
		s.removeInstances(i40, i40Deployed[:i40DeployedCount]), err
}

func (s *BestFitStrategy) fillMachines(machines []*cloud.Machine, instances []*cloud.Instance) (restInstances []*cloud.Instance, err error) {
	//s.analysisMachineDiskDistribution(machines)
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
	for i, m := range machines {
		if i > 0 && i%100 == 0 {
			fmt.Println("fillMachines", i)
		}

		i80Rest, i60Rest, i40Rest, err = s.fillMachine(m, i80Rest, i60Rest, i40Rest)
		if err != nil {
			fmt.Println("fillMachines", i)
			return nil, err
		}

		if m.LevelConfig.Disk == 1024 {
			if m.Disk < 1000 {
				return nil, fmt.Errorf("not filled,disk=%d,%d", m.Disk, i)
			}
		}

		if len(i60Rest) == 0 {
			fmt.Println("great finish")
			break
		}

		//fmt.Println("fillMachines", len(i80Rest), len(i60Rest), len(i40Rest))

		n := 0
		for _, mm := range machines {
			n += mm.InstanceArrayCount
		}

		nn := len(i80Rest) + len(i60Rest) + len(i40Rest)
		if nn+n != 68219 {
			panic(fmt.Errorf("n=%d,nn=%d,80=%d,60=%d,40=%d\n", n, nn, len(i80Rest), len(i60Rest), len(i40Rest)))
		}
	}

	cloud.SetDebug(false)

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
					m.AddInstance(instance)
					deployedInstances[instance.InstanceId] = instance
					deployed167 = true
					break
				}
			}
			if !deployed167 {
				machineList[machineIndex].AddInstance(instance)
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
				m.AddInstance(instance)
				deployedInstances[instance.InstanceId] = instance
				deployed100 = true
				break
			}
			if !deployed100 {
				machineList[machineIndex].AddInstance(instance)
				deployedInstances[instance.InstanceId] = instance
				deployedMachines = append(deployedMachines, machineList[machineIndex])
				machineIndex++
			}
		} else {
			machineList[machineIndex].AddInstance(instance)
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

				m.AddInstance(instance)
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
				m.AddInstance(instance)
				deployedInstances[instance.InstanceId] = instance
				deployed150 = true
				break
			}
		}
		if !deployed150 {
			machineList[machineIndex].AddInstance(instance)
			deployedInstances[instance.InstanceId] = instance
			deployedMachines = append(deployedMachines, machineList[machineIndex])
			machineIndex++
		}
	}

	for ; machineIndex < 5505; machineIndex++ {
		deployedMachines = append(deployedMachines, machineList[machineIndex])
	}

	//填充这些机器
	restInstances = s.preDeployRemoveDeployed(instanceList, deployedInstances)
	restInstances, err = s.fillMachines(deployedMachines, restInstances)
	if err != nil {
		//fmt.Println("lalala")
		for _, v := range machineList {
			if v.Disk > 0 {
				if v.LevelConfig.Disk == 1024 {
					if v.Disk <= 1000 {
						//fmt.Println(i, v.Disk)
					}
				} else {
					if v.Disk < 600 {
						//fmt.Println(i, v.Disk)
					}
				}
			}
		}
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
