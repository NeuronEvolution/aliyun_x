package main

import (
	"bytes"
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"github.com/NeuronEvolution/aliyun_x/strategies/bfs"
	"io/ioutil"
	"os"
	"strconv"
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

	instanceMachineList = instanceMachineList[:0]
	//instanceList = instanceList[:20000]
	begin := time.Now()
	r := cloud.NewResourceManagement()
	r.SetStrategy(bfs.NewFreeSmallerStrategy(r))
	err = r.Init(machineResourceDataList, appResourcesDataList, appInterferenceDataList, instanceMachineList)
	if err != nil {
		fmt.Printf("r.Init failed,%s", err)
		return
	}

	r.DebugPrintStatus()

	err = r.PostInit()
	if err != nil {
		fmt.Printf("r.PostInit failed,%s", err)
		return
	}

	err = r.AddInstanceList(instanceDeployDataList)
	if err != nil {
		fmt.Println(err)
		return
	}

	end := time.Now()

	r.DebugPrintStatus()

	fmt.Printf("\n\n\n")
	fmt.Printf("*****************************************************************\n")
	fmt.Printf("*****************************************************************\n")
	fmt.Printf("*****************************************************************\n")

	playback := cloud.NewResourceManagement()
	err = playback.Init(machineResourceDataList, appResourcesDataList, appInterferenceDataList, instanceMachineList)
	if err != nil {
		fmt.Printf("r.Init failed,%s", err)
		return
	}
	playback.DebugPrintStatus()

	err = playback.Play(r.DeployCommandHistory)
	if err != nil {
		fmt.Printf("playback.Play failed,%s", err)
		return
	}

	playback.DebugPrintStatus()

	fmt.Printf("time=%f\n", end.Sub(begin).Seconds())

	outputFile := fmt.Sprintf("_output/submit_%s", time.Now().Format("20060102_150405"))
	buf := bytes.NewBufferString("")
	for _, v := range r.DeployCommandHistory.List[:r.DeployCommandHistory.ListCount] {
		buf.WriteString("inst_")
		buf.WriteString(strconv.Itoa(v.InstanceId))
		buf.WriteString(",")
		buf.WriteString("machine_")
		buf.WriteString(strconv.Itoa(v.MachineId))
		buf.WriteString("\n")
	}

	err = ioutil.WriteFile(outputFile+".csv", buf.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	summaryBuf := bytes.NewBufferString("")
	r.DebugStatus(summaryBuf)
	summaryBuf.WriteString(fmt.Sprintf("time=%f\n", end.Sub(begin).Seconds()))
	err = ioutil.WriteFile(fmt.Sprintf(outputFile+"_summary.csv"),
		summaryBuf.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}
