package main

import (
	"encoding/csv"
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"os"
	"strconv"
	"strings"
)

func trimApp(s string) string {
	return s[4:]
}

func trimMachine(s string) string {
	return s[8:]
}

func trimInstance(s string) string {
	return s[5:]
}

func loadCsv(file string) (data [][]string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	data, err = r.ReadAll()
	if err != nil {
		return nil, err
	}

	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("no data")
	}

	return data, nil
}

func loadAppInterferenceData(file string) (result []*cloud.AppInterferenceConfig, err error) {
	data, err := loadCsv(file)
	if err != nil {
		return nil, err
	}

	result = make([]*cloud.AppInterferenceConfig, len(data))
	for i, v := range data {
		if len(v) < 3 {
			return nil, fmt.Errorf("loadAppInterference data row len<3")
		}

		item := &cloud.AppInterferenceConfig{}
		column := 0

		item.AppId1, err = strconv.Atoi(trimApp(v[column]))
		if err != nil {
			return nil, err
		}
		column++

		item.AppId2, err = strconv.Atoi(trimApp(v[column]))
		if err != nil {
			return nil, err
		}
		column++

		item.Interference, err = strconv.Atoi(v[column])
		if err != nil {
			return nil, err
		}
		column++

		result[i] = item
	}
	return result, nil
}

func loadAppResourceData(file string) (result []*cloud.AppResourcesConfig, err error) {
	data, err := loadCsv(file)
	if err != nil {
		return nil, err
	}

	result = make([]*cloud.AppResourcesConfig, len(data))
	for i, v := range data {
		if len(v) < 6 {
			return nil, fmt.Errorf("loadAppResource data row len<6")
		}

		item := &cloud.AppResourcesConfig{}
		column := 0

		item.AppId, err = strconv.Atoi(trimApp(v[column]))
		if err != nil {
			return nil, err
		}
		column++

		cpuTokens := strings.Split(v[column], "|")
		if len(cpuTokens) != len(item.Cpu) {
			return nil,
				fmt.Errorf("loadAppResource cpu len %d failed,required %d", len(cpuTokens), len(item.Cpu))
		}
		for tokenIndex, token := range cpuTokens {
			item.Cpu[tokenIndex], err = strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, err
			}
		}
		column++

		memTokens := strings.Split(v[column], "|")
		if len(memTokens) != len(item.Mem) {
			return nil,
				fmt.Errorf("loadAppResource mem len %d failed,required %d", len(cpuTokens), len(item.Cpu))
		}
		for tokenIndex, token := range memTokens {
			item.Mem[tokenIndex], err = strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, err
			}
		}
		column++

		item.Disk, err = strconv.Atoi(v[column])
		if err != nil {
			return nil, err
		}
		column++

		item.P, err = strconv.Atoi(v[column])
		if err != nil {
			return nil, err
		}
		column++

		item.M, err = strconv.Atoi(v[column])
		if err != nil {
			return nil, err
		}
		column++

		item.PM, err = strconv.Atoi(v[column])
		if err != nil {
			return nil, err
		}
		column++

		result[i] = item
	}

	return result, nil
}

func loadInstanceDeployData(file string) (result []*cloud.InstanceDeployConfig, err error) {
	data, err := loadCsv(file)
	if err != nil {
		return nil, err
	}

	result = make([]*cloud.InstanceDeployConfig, len(data))
	for i, v := range data {
		item := &cloud.InstanceDeployConfig{}
		column := 0

		item.InstanceId, err = strconv.Atoi(trimInstance(v[column]))
		if err != nil {
			return nil, err
		}
		column++

		item.AppId, err = strconv.Atoi(trimApp(v[column]))
		if err != nil {
			return nil, err
		}
		column++

		if v[column] != "" {
			item.MachineId, err = strconv.Atoi(trimMachine(v[column]))
			if err != nil {
				return nil, err
			}
		}

		column++

		result[i] = item
	}

	return result, nil
}

func loadMachineResourcesData(file string) (result []*cloud.MachineResourcesConfig, err error) {
	data, err := loadCsv(file)
	if err != nil {
		return nil, err
	}

	result = make([]*cloud.MachineResourcesConfig, len(data))
	for i, v := range data {
		item := &cloud.MachineResourcesConfig{}
		column := 0

		item.MachineId, err = strconv.Atoi(trimMachine(v[column]))
		if err != nil {
			return nil, err
		}
		column++

		item.Cpu, err = strconv.ParseFloat(v[column], 64)
		if err != nil {
			return nil, err
		}
		column++

		item.Mem, err = strconv.ParseFloat(v[column], 64)
		if err != nil {
			return nil, err
		}
		column++

		item.Disk, err = strconv.Atoi(v[column])
		if err != nil {
			return nil, err
		}
		column++

		item.P, err = strconv.Atoi(v[column])
		if err != nil {
			return nil, err
		}
		column++

		item.M, err = strconv.Atoi(v[column])
		if err != nil {
			return nil, err
		}
		column++

		item.PM, err = strconv.Atoi(v[column])
		if err != nil {
			return nil, err
		}
		column++

		result[i] = item
	}

	return result, nil
}
