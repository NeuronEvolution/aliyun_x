package fss

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"time"
)

func (s *FreeSmallerStrategy) getAppInferenceMachineList() []*cloud.Machine {
	machineList := make([]*cloud.Machine, 0)
	for _, level := range s.R.MachineDeployPool.MachineLevelDeployArray {
		for _, v := range level.MachineCollection.List[:level.MachineCollection.ListCount] {
			if v.HasBadConstraint() {
				machineList = append(machineList, v)
			}
		}
	}

	return machineList
}

func (s *FreeSmallerStrategy) resolveAppInference() (err error) {
	fmt.Printf("FreeSmallerStrategy.resolveAppInference\n")
	machineList := s.getAppInferenceMachineList()
	instanceList := make([]*cloud.Instance, 0)
	for _, m := range machineList {
		instanceList = append(instanceList, m.InstanceArray[:m.InstanceArrayCount]...)
	}
	fmt.Printf("FreeSmallerStrategy.resolveAppInference begin machineCount=%d,instanceCount=%d\n",
		len(machineList), len(instanceList))

	err = s.redeployInstanceList(instanceList, true)
	if err != nil {
		return err
	}

	machineList = s.getAppInferenceMachineList()
	if len(machineList) != 0 {
		return fmt.Errorf("FreeSmallerStrategy.resolveAppInference not completed")
	}

	return nil
}

func (s *FreeSmallerStrategy) getHighCpuMachineList() []*cloud.Machine {
	machineList := make([]*cloud.Machine, 0)
	for _, level := range s.R.MachineDeployPool.MachineLevelDeployArray {
		for _, v := range level.MachineCollection.List[:level.MachineCollection.ListCount] {
			if (v.LevelConfig.Cpu == cloud.MachineCpuMax && v.GetCost() > HighLevelCpuMax) ||
				(v.LevelConfig.Cpu < cloud.MachineCpuMax && v.GetCost() > LowLevelCpuMax) {
				machineList = append(machineList, v)
			}
		}
	}

	return machineList
}

func (s *FreeSmallerStrategy) getLowLevelMachineList() []*cloud.Machine {
	machineList := make([]*cloud.Machine, 0)

	if s.R.MachineDeployPool.MachineLevelDeployArray.Len() == 0 {
		return machineList
	}

	if s.R.MachineFreePool.MachineLevelFreeArray.Len() <= 1 {
		return machineList
	}

	startLevel := 0
	if s.R.MachineDeployPool.MachineLevelDeployArray[0].LevelConfig.Cpu ==
		s.R.MachineFreePool.MachineLevelFreeArray[0].LevelConfig.Cpu {
		startLevel = 1
	}
	for _, level := range s.R.MachineDeployPool.MachineLevelDeployArray[startLevel:] {
		if level.MachineCollection.ListCount > 0 {
			machineList = append(machineList, level.MachineCollection.List[:level.MachineCollection.ListCount]...)
		}
	}

	return machineList
}

func (s *FreeSmallerStrategy) resolveHighCpuAndLowLevelMachine() (err error) {
	fmt.Printf("FreeSmallerStrategy.resolveHighCpuAndLowLevelMachine\n")

	highCpuMachineList := s.getHighCpuMachineList()
	fmt.Printf("FreeSmallerStrategy.resolveHighCpuAndLowLevelMachine begin hight cpu machine count=%d\n",
		len(highCpuMachineList))
	lowLevelMachineList := s.getLowLevelMachineList()
	fmt.Printf("FreeSmallerStrategy.resolveHighCpuAndLowLevelMachine begin low level machine count=%d\n",
		len(lowLevelMachineList))

	instanceList := make([]*cloud.Instance, 0)
	for _, m := range highCpuMachineList {
		instanceList = append(instanceList, m.InstanceArray[:m.InstanceArrayCount]...)
	}
	for _, m := range lowLevelMachineList {
		instanceList = append(instanceList, m.InstanceArray[:m.InstanceArrayCount]...)
	}

	fmt.Printf("FreeSmallerStrategy.resolveHighCpuAndLowLevelMachine instance count=%d\n",
		len(instanceList))

	err = s.redeployInstanceList(instanceList, false)
	if err != nil {
		return nil
	}

	highCpuMachineList = s.getHighCpuMachineList()
	fmt.Printf("FreeSmallerStrategy.resolveHighCpuAndLowLevelMachine end hight cpu machine count=%d\n",
		len(highCpuMachineList))
	lowLevelMachineList = s.getLowLevelMachineList()
	fmt.Printf("FreeSmallerStrategy.resolveHighCpuAndLowLevelMachine end low level machine count=%d\n",
		len(lowLevelMachineList))

	return nil
}

func (s *FreeSmallerStrategy) PostInit() (err error) {
	begin := time.Now()
	err = s.resolveAppInference()
	if err != nil {
		return err
	}
	fmt.Printf("FreeSmallerStrategy.PostInit resolveAppInference time=%f\n", time.Now().Sub(begin).Seconds())
	s.R.DebugPrintStatus()

	begin = time.Now()
	err = s.resolveHighCpuAndLowLevelMachine()
	if err != nil {
		return err
	}
	fmt.Printf("FreeSmallerStrategy.PostInit resolveHighCpuAndLowLevelMachine time=%f\n", time.Now().Sub(begin).Seconds())
	s.R.DebugPrintStatus()

	return nil
}
