package bfs_v2

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"os"
	"time"
)

func (s *Strategy) mergeFinal() {
	startCost := s.R.CalculateTotalCostScore()
	fmt.Println("mergeFinal start cpu cost", startCost)

	lastSaveTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)

	pTable := s.randBig2Big(len(s.machineDeployList))

	currentCost := startCost
	loop := 0
	deadLoop := 0
	for ; ; loop++ {
		cloud.SortMachineByCpuCost(s.machineDeployList)
		machinesByCpu := s.randMachinesBig2Big(s.machineDeployList, pTable, 32)
		ok, delta := s.BatchBestMergeMachines(machinesByCpu, deadLoop)
		if !ok {
			fmt.Println("mergeFinal dead loop", deadLoop)
			deadLoop++
			if deadLoop > 256 {
				break
			}

			continue
		}

		deadLoop = 0
		currentCost += delta
		fmt.Printf("mergeFinal loop %d %f %f\n", loop, startCost, currentCost)

		now := time.Now()
		if now.Sub(lastSaveTime).Seconds() > 60*30 {
			fmt.Printf("mergeFinal save loop %d %f %f\n", loop, startCost, currentCost)
			err := s.R.MergeAndOutput()
			if err != nil {
				fmt.Println("mergeFinal save failed", err)
			}
			lastSaveTime = now
		}

		_, err := os.Stat("aliyun_stop")
		if err == nil {
			fmt.Println("mergeFinal aliyun_stop")
			break
		}
	}

	fmt.Printf("mergeFinal end %d %f %f\n", loop, startCost, s.R.CalculateTotalCostScore())
}
