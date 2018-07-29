package main

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"github.com/NeuronEvolution/aliyun_x/strategies/bfs_v2"
	"time"
)

//const appInterferenceFile = "./data/scheduling_preliminary_app_interference_20180606.csv"
//const appResourcesFile = "./data/scheduling_preliminary_app_resources_20180606.csv"
//const instanceDeployFile = "./data/scheduling_preliminary_instance_deploy_20180606.csv"
//const machineResourceFile = "./data/scheduling_preliminary_machine_resources_20180606.csv"

const appInterferenceFile = "./data/scheduling_preliminary_b_app_interference_20180726.csv"
const appResourcesFile = "./data/scheduling_preliminary_b_app_resources_20180726.csv"
const instanceDeployFile = "./data/scheduling_preliminary_b_instance_deploy_20180726.csv"
const machineResourceFile = "./data/scheduling_preliminary_b_machine_resources_20180726.csv"

func testCase(machines []*cloud.MachineResourcesConfig, instances []*cloud.InstanceDeployConfig,
	highMachineCount int, lowMachineCount int, instanceCount int) (
	restMachines []*cloud.MachineResourcesConfig, restInstances []*cloud.InstanceDeployConfig) {

	highCount := 0
	lowCount := 0
	for _, m := range machines {
		if m.MachineLevelConfig.Disk == cloud.HighDisk {
			restMachines = append(restMachines, m)
			highCount++
		} else if m.MachineLevelConfig.Disk == cloud.LowDisk {
			restMachines = append(restMachines, m)
			lowCount++
		}
	}

	return nil, nil
}

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

	fmt.Printf("DataSize\n")
	fmt.Printf("   appInterferenceDataList=%d\n", len(appInterferenceDataList))
	fmt.Printf("   appResourcesDataList=%d\n", len(appResourcesDataList))
	fmt.Printf("   instanceDeployDataList=%d\n", len(instanceDeployDataList))
	fmt.Printf("   machineResourceDataList=%d\n", len(machineResourceDataList))

	for _, v := range appResourcesDataList {
		v.CalcTimedResourceStatistics()

		for _, inference := range appInterferenceDataList {
			if inference.AppId1 == v.AppId || inference.AppId2 == v.AppId {
				v.InferenceAppCount++
			}
		}
	}

	//数据分析
	analysis := NewAnalysisContext(appInterferenceDataList, appResourcesDataList, machineResourceDataList, instanceDeployDataList)
	analysis.Run()

	//调度
	begin := time.Now()
	result, err := cloud.Run(
		machineResourceDataList,
		appResourcesDataList,
		appInterferenceDataList,
		instanceDeployDataList, func(r *cloud.ResourceManagement) cloud.Strategy {
			return bfs_v2.NewStrategy(r)
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	end := time.Now()

	//输出
	err = output(result, end.Sub(begin))
	if err != nil {
		fmt.Println(err)
	}
}
