package main

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"github.com/NeuronEvolution/aliyun_x/strategies/sffs"
	"time"
)

const appInterferenceFile = "./data/scheduling_preliminary_app_interference_20180606.csv"
const appResourcesFile = "./data/scheduling_preliminary_app_resources_20180606.csv"
const instanceDeployFile = "./data/scheduling_preliminary_instance_deploy_20180606.csv"
const machineResourceFile = "./data/scheduling_preliminary_machine_resources_20180606.csv"

func main() {
	appInterferenceDataList, err := loadAppInterferenceData(appInterferenceFile)
	if err != nil {
		fmt.Println("loadAppInterferenceData failed", err)
		return
	}

	appResourcesDataList, err := loadAppResourceData(appResourcesFile)
	if err != nil {
		fmt.Println("loadAppResourceData failed", err)
		return
	}

	instanceDeployDataList, err := loadInstanceDeployData(instanceDeployFile)
	if err != nil {
		fmt.Println("loadInstanceDeployData failed", err)
		return
	}

	machineResourceDataList, err := loadMachineResourcesData(machineResourceFile)
	if err != nil {
		fmt.Println("loadMachineResourcesData failed", err)
		return
	}

	fmt.Println("appInterferenceDataList", len(appInterferenceDataList))
	fmt.Println("appResourcesDataList", len(appResourcesDataList))
	fmt.Println("instanceDeployDataList", len(instanceDeployDataList))
	fmt.Println("machineResourceDataList", len(machineResourceDataList))

	//导入机器数据
	r := cloud.NewResourceManagement()
	r.SetStrategy(sffs.NewASortedFirstFitStrategy(r))
	//r.SetStrategy(strategies.NewFreeSmallerStrategy(r))

	var maxMachineId, maxAppId, maxInstanceId int

	for _, v := range machineResourceDataList {
		if v.MachineId > maxMachineId {
			maxMachineId = v.MachineId
		}
	}

	//导入应用资源数据
	for _, v := range appResourcesDataList {
		if v.AppId > maxAppId {
			maxAppId = v.AppId
		}
	}

	//导入应用冲突数据
	for _, v := range appInterferenceDataList {
		if v.AppId1 > maxAppId {
			maxAppId = v.AppId1
		}
		if v.AppId2 > maxAppId {
			maxAppId = v.AppId2
		}
	}

	instanceList := make([]*cloud.InstanceDeployConfig, 0)
	instanceMachineList := make([]*cloud.InstanceDeployConfig, 0)
	for _, v := range instanceDeployDataList {
		if v.MachineId == 0 {
			instanceList = append(instanceList, v)
		} else {
			instanceMachineList = append(instanceMachineList, v)
		}
		if v.AppId > maxAppId {
			maxAppId = v.AppId
		}
		if v.MachineId > maxMachineId {
			maxMachineId = v.MachineId
		}
		if v.InstanceId > maxInstanceId {
			maxInstanceId = v.InstanceId
		}
	}

	fmt.Printf("maxAppId=%d,maxInstanceId=%d,maxMachineId=%d\n", maxAppId, maxInstanceId, maxMachineId)
	fmt.Printf("deployed=%d,non-deployed=%d,total=%d\n",
		len(instanceMachineList), len(instanceList), len(instanceDeployDataList))

	begin := time.Now()

	err = r.Init(machineResourceDataList, appResourcesDataList, appInterferenceDataList, instanceMachineList)
	if err != nil {
		fmt.Printf("r.Init failed,%s", err)
		return
	}

	cloud.SetDebug(true)
	r.DebugDeployStatus()
	err = r.ResolveAppInference()
	if err != nil {
		fmt.Printf("r.ResolveAppInference failed,%s", err)
	}

	r.BatchAddInstance(instanceList)
	end := time.Now()

	r.DebugDeployStatus()

	playback := cloud.NewResourceManagement()
	err = playback.Init(machineResourceDataList, appResourcesDataList, appInterferenceDataList, instanceMachineList)
	if err != nil {
		fmt.Printf("r.Init failed,%s", err)
		return
	}

	playback.DebugDeployStatus()

	err = playback.Play(r.DeployCommandHistory)
	if err != nil {
		fmt.Printf("playback.Play failed,%s", err)
		return
	}

	playback.DebugDeployStatus()

	fmt.Printf("time=%10f\n", end.Sub(begin).Seconds())
}
