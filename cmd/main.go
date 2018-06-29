package main

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"github.com/NeuronEvolution/aliyun_x/strategies"
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
	r.SetStrategy(strategies.NewAllocMachineIfDeployFailedStrategy(r))
	//r.SetStrategy(strategies.NewFreeSmallerStrategy(r))

	for _, v := range machineResourceDataList {
		r.AddMachine(v)
	}

	//导入应用资源数据
	for _, v := range appResourcesDataList {
		r.SaveAppResourceConfig(v)
	}

	//导入应用冲突数据
	for _, v := range appInterferenceDataList {
		r.SaveAppInterferenceConfig(v)
	}

	instanceList := make([]*cloud.InstanceDeployConfig, 0)
	instanceMachineList := make([]*cloud.InstanceDeployConfig, 0)
	for _, v := range instanceDeployDataList {
		if v.MachineId == 0 {
			instanceList = append(instanceList, v)
		} else {
			instanceMachineList = append(instanceMachineList, v)
		}
	}

	begin := time.Now()

	cloud.SetDebug(true)
	//r.InitInstanceDeploy(instanceMachineList)

	r.ResetCommandHistory()

	r.BatchAddInstance(instanceDeployDataList)
	end := time.Now()

	r.DebugDeployStatus()
	fmt.Printf("%d %d %d\n", len(instanceMachineList), len(instanceList), len(instanceDeployDataList))
	fmt.Printf("time %10.2f\n", end.Sub(begin).Minutes())
	fmt.Println("cost", r.CalculateTotalCostScore())
}
