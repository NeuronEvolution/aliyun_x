package fss

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"sort"
	"time"
)

func (s *FreeSmallerStrategy) resolveAppInference() (err error) {
	fmt.Printf("FreeSmallerStrategy.resolveAppInference\n")
	for i := 0; ; i++ {
		//fmt.Printf("FreeSmallerStrategy.resolveAppInference %d\n", i)
		var m *cloud.Machine
		for _, level := range s.R.MachineDeployPool.MachineLevelDeployArray {
			for _, v := range level.MachineCollection.List[:level.MachineCollection.ListCount] {
				if v.HasBadConstraint() {
					m = v
					break
				}
			}
			if m != nil {
				break
			}
		}

		if m == nil {
			fmt.Printf("SortedFirstFitStrategy.resolveAppInference total expand %d\n", i)
			break
		}

		err := s.redeployMachine(m, false)
		if err != nil {
			return err
		}
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
	sort.Sort(cloud.MachineListSortByCostDesc(machineList))

	return machineList
}

func (s *FreeSmallerStrategy) resolveHighCpu() error {
	fmt.Printf("FreeSmallerStrategy.resolveHighCpu\n")
	machineRedeployList := s.getHighCpuMachineList()
	fmt.Printf("FreeSmallerStrategy.resolveHighCpu high cpu before count=%d\n", len(machineRedeployList))

	for _, v := range machineRedeployList {
		s.redeployMachine(v, false)
	}

	machineRedeployList = s.getHighCpuMachineList()
	fmt.Printf("FreeSmallerStrategy.resolveHighCpu high cpu after count=%d\n", len(machineRedeployList))
	for _, v := range machineRedeployList {
		fmt.Printf("    %f,%d\n", v.GetCost(), v.MachineId)
		//v.DebugPrint()
	}
	fmt.Println("")

	return nil
}

func (s *FreeSmallerStrategy) resolveSmallMachine() error {
	fmt.Printf("FreeSmallerStrategy.resolveSmallMachine\n")
	if s.R.MachineDeployPool.MachineLevelDeployArray.Len() == 0 {
		return nil
	}

	if s.R.MachineFreePool.MachineLevelFreeArray.Len() <= 1 {
		return nil
	}

	countResolved := 0
	for {
		hasHighLevelMachine := false
		for _, v := range s.R.MachineFreePool.MachineLevelFreeArray[:s.R.MachineFreePool.MachineLevelFreeArray.Len()-1] {
			if v.MachineCollection.ListCount > 0 {
				hasHighLevelMachine = true
				break
			}
		}
		if !hasHighLevelMachine {
			break
		}

		startLevel := 0
		if s.R.MachineDeployPool.MachineLevelDeployArray[0].LevelConfig.Cpu ==
			s.R.MachineFreePool.MachineLevelFreeArray[0].LevelConfig.Cpu {
			startLevel = 1
		}

		hasLowMachine := false
		for _, level := range s.R.MachineDeployPool.MachineLevelDeployArray[startLevel:] {
			if level.MachineCollection.ListCount > 0 {
				hasLowMachine = true
				s.redeployMachine(level.MachineCollection.List[0], false)
				countResolved++
				break
			}
		}

		if !hasLowMachine {
			break
		}
	}

	fmt.Printf("FreeSmallerStrategy.resolveSmallMachine countResolved=%d\n", countResolved)

	return nil
}

func (s *FreeSmallerStrategy) PostInit() (err error) {
	begin := time.Now()
	err = s.resolveAppInference()
	if err != nil {
		return err
	}
	s.R.DebugPrintStatus()
	fmt.Printf("FreeSmallerStrategy.PostInit resolveAppInference time=%f\n", time.Now().Sub(begin).Seconds())

	begin = time.Now()
	err = s.resolveHighCpu()
	if err != nil {
		return err
	}
	s.R.DebugPrintStatus()
	fmt.Printf("FreeSmallerStrategy.PostInit resolveHighCpu time=%f\n", time.Now().Sub(begin).Seconds())

	begin = time.Now()
	err = s.resolveSmallMachine()
	if err != nil {
		return err
	}
	s.R.DebugPrintStatus()
	fmt.Printf("FreeSmallerStrategy.PostInit resolveSmallMachine time=%f\n", time.Now().Sub(begin).Seconds())

	return nil
}
