package bfs

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"time"
)

func (s *BestFitStrategy) getAppInferenceMachineList() []*cloud.Machine {
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

func (s *BestFitStrategy) resolveAppInference() (err error) {
	fmt.Printf("BestFitStrategy.resolveAppInference\n")
	machineList := s.getAppInferenceMachineList()
	instanceList := make([]*cloud.Instance, 0)
	for _, m := range machineList {
		instanceList = append(instanceList, m.InstanceArray[:m.InstanceArrayCount]...)
	}
	fmt.Printf("BestFitStrategy.resolveAppInference begin machineCount=%d,instanceCount=%d\n",
		len(machineList), len(instanceList))

	err = s.redeployInstanceList(instanceList, true)
	if err != nil {
		return err
	}

	machineList = s.getAppInferenceMachineList()
	if len(machineList) != 0 {
		return fmt.Errorf("BestFitStrategy.resolveAppInference not completed")
	}

	return nil
}

func (s *BestFitStrategy) getHighCpuMachineList() []*cloud.Machine {
	machineList := make([]*cloud.Machine, 0)
	for _, level := range s.R.MachineDeployPool.MachineLevelDeployArray {
		for _, v := range level.MachineCollection.List[:level.MachineCollection.ListCount] {
			if v.GetCost() > cloud.MaxCpu {
				machineList = append(machineList, v)
			}
		}
	}

	return machineList
}

func (s *BestFitStrategy) getLowLevelMachineList() []*cloud.Machine {
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

func (s *BestFitStrategy) resolveHighCpuMachine() (err error) {
	fmt.Printf("BestFitStrategy.resolveHighCpuMachine\n")

	highCpuMachineList := s.getHighCpuMachineList()
	instanceList := make([]*cloud.Instance, 0)
	for _, m := range highCpuMachineList {
		instanceList = append(instanceList, m.InstanceArray[:m.InstanceArrayCount]...)
	}

	fmt.Printf("BestFitStrategy.resolveHighCpuMachine begin hight cpu machineCount=%d,instanceCount=%d\n",
		len(highCpuMachineList), len(instanceList))

	err = s.redeployInstanceList(instanceList, false)
	if err != nil {
		return nil
	}

	highCpuMachineList = s.getHighCpuMachineList()
	fmt.Printf("BestFitStrategy.resolveHighCpuMachine end hight cpu machine count=%d\n",
		len(highCpuMachineList))

	return nil
}

func (s *BestFitStrategy) resolveLowLevelMachine() (err error) {
	fmt.Printf("BestFitStrategy.resolveLowLevelMachine\n")

	lowLevelMachineList := s.getLowLevelMachineList()
	instanceList := make([]*cloud.Instance, 0)
	for _, m := range lowLevelMachineList {
		instanceList = append(instanceList, m.InstanceArray[:m.InstanceArrayCount]...)
	}

	fmt.Printf("BestFitStrategy.resolveLowLevelMachine begin low level machineCount=%d,instanceCount=%d\n",
		len(lowLevelMachineList), len(instanceList))

	err = s.redeployInstanceList(instanceList, false)
	if err != nil {
		return nil
	}

	lowLevelMachineList = s.getLowLevelMachineList()
	fmt.Printf("BestFitStrategy.resolveLowLevelMachine end low level machine count=%d\n",
		len(lowLevelMachineList))

	return nil
}

func (s *BestFitStrategy) PostInit() (err error) {
	s.machineDeployList = s.getDeployMachineList(MachineDeployCount)
	if len(s.machineDeployList) != MachineDeployCount {
		panic("BestFitStrategy.AddInstanceList getDeployMachineList failed")
	}

	begin := time.Now()
	err = s.resolveAppInference()
	if err != nil {
		return err
	}
	fmt.Printf("BestFitStrategy.PostInit resolveAppInference time=%f\n", time.Now().Sub(begin).Seconds())
	s.R.DebugPrintStatus()

	begin = time.Now()
	//err = s.resolveHighCpuMachine()
	if err != nil {
		return err
	}
	fmt.Printf("BestFitStrategy.PostInit resolveHighCpuMachine time=%f\n", time.Now().Sub(begin).Seconds())
	s.R.DebugPrintStatus()

	begin = time.Now()
	//err = s.resolveLowLevelMachine()
	if err != nil {
		return err
	}
	fmt.Printf("BestFitStrategy.PostInit resolveLowLevelMachine time=%f\n", time.Now().Sub(begin).Seconds())
	s.R.DebugPrintStatus()

	return nil
}
