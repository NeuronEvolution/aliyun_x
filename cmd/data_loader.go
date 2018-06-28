package main

import (
	"encoding/csv"
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/clound"
	"os"
	"strconv"
	"strings"
)

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

func loadAppInterferenceData(file string) (result []*clound.AppInterferenceConfig, err error) {
	data, err := loadCsv(file)
	if err != nil {
		return nil, err
	}

	result = make([]*clound.AppInterferenceConfig, len(data))
	for i, v := range data {
		if len(v) < 3 {
			return nil, fmt.Errorf("loadAppInterference data row len<3")
		}

		item := &clound.AppInterferenceConfig{}
		column := 0

		item.AppId1 = v[column]
		column++

		item.AppId2 = v[column]
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

func loadAppResourceData(file string) (result []*clound.AppResourcesConfig, err error) {
	data, err := loadCsv(file)
	if err != nil {
		return nil, err
	}

	result = make([]*clound.AppResourcesConfig, len(data))
	for i, v := range data {
		if len(v) < 6 {
			return nil, fmt.Errorf("loadAppResource data row len<6")
		}

		item := &clound.AppResourcesConfig{}
		column := 0

		item.AppId = v[column]
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

func loadInstanceDeployData(file string) (result []*clound.InstanceDeployConfig, err error) {
	data, err := loadCsv(file)
	if err != nil {
		return nil, err
	}

	result = make([]*clound.InstanceDeployConfig, len(data))
	for i, v := range data {
		item := &clound.InstanceDeployConfig{}
		column := 0

		item.InstanceId = v[column]
		column++

		item.AppId = v[column]
		column++

		item.MachineId = v[column]
		column++

		result[i] = item
	}

	return result, nil
}

func loadMachineResourcesData(file string) (result []*clound.MachineResourcesConfig, err error) {
	data, err := loadCsv(file)
	if err != nil {
		return nil, err
	}

	result = make([]*clound.MachineResourcesConfig, len(data))
	for i, v := range data {
		item := &clound.MachineResourcesConfig{}
		column := 0

		item.MachineId = v[column]
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

		item.P, err = strconv.Atoi(v[column])
		if err != nil {
			return nil, err
		}
		column++

		result[i] = item
	}

	return result, nil
}
