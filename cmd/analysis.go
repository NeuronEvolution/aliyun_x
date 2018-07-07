package main

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

func analysis(
	appInterferenceList []*cloud.AppInterferenceConfig,
	appResourcesList []*cloud.AppResourcesConfig,
	machineResourcesList []*cloud.MachineResourcesConfig,
	instanceDeployList []*cloud.InstanceDeployConfig) {

	diskDist := make(map[int]int, 0)
	for _, v := range instanceDeployList {
		for _, a := range appResourcesList {
			if a.AppId == v.AppId {
				_, has := diskDist[a.Disk]
				if !has {
					diskDist[a.Disk] = 1
				} else {
					diskDist[a.Disk]++
				}
				break
			}
		}
	}

	for disk, count := range diskDist {
		fmt.Printf("disk=%d,count=%d\n", disk, count)
	}
}
