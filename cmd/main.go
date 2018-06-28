package main

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/clound"
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
	r := clound.NewResourceManagement()
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

	clound.SetDebug(true)

	//一个一个部署
	for i, v := range instanceDeployDataList {
		err = r.AddInstance(v.InstanceId, v.AppId)
		if err != nil {
			fmt.Printf("AddInstance %d failed %s", i, err)
			break
		}

		if i >= 100 {
			break
		}
	}
}
