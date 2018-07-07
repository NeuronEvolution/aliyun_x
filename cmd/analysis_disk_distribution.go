package main

import (
	"fmt"
	"sort"
)

type diskCount struct {
	disk  int
	count int
}

func (c *AnalysisContext) AnalysisDiskDistribution() {
	appDiskDist := make(map[int]int, 0)
	for _, a := range c.appResourcesList {
		_, has := appDiskDist[a.Disk]
		if !has {
			appDiskDist[a.Disk] = 1
		} else {
			appDiskDist[a.Disk]++
		}
	}
	appDiskCountList := make([]*diskCount, 0)
	for i, v := range appDiskDist {
		appDiskCountList = append(appDiskCountList, &diskCount{disk: i, count: v})
	}
	sort.Slice(appDiskCountList, func(i, j int) bool {
		return appDiskCountList[i].disk > appDiskCountList[j].disk
	})
	fmt.Printf("AnalysisDiskDistribution by app\n")
	for _, v := range appDiskCountList {
		fmt.Printf("    disk=%d,count=%d\n", v.disk, v.count)
	}

	instanceDiskDist := make(map[int]int, 0)
	for _, v := range c.instanceDeployList {
		appResource := c.appResourcesMap[v.AppId]
		_, has := instanceDiskDist[appResource.Disk]
		if !has {
			instanceDiskDist[appResource.Disk] = 1
		} else {
			instanceDiskDist[appResource.Disk]++
		}
	}
	instanceDiskCountList := make([]*diskCount, 0)
	for i, v := range appDiskDist {
		instanceDiskCountList = append(instanceDiskCountList, &diskCount{disk: i, count: v})
	}
	sort.Slice(instanceDiskCountList, func(i, j int) bool {
		return instanceDiskCountList[i].disk > instanceDiskCountList[j].disk
	})
	fmt.Printf("AnalysisDiskDistribution by instance\n")
	for _, v := range instanceDiskCountList {
		fmt.Printf("    disk=%d,count=%d\n", v.disk, v.count)
	}
}
