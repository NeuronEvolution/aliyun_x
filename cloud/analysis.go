package cloud

import (
	"fmt"
	"sort"
)

type DiskCount struct {
	Disk  int
	Count int
}

func AnalysisDiskDistributionByInstance(instanceList []*Instance) {
	fmt.Printf("AnalysisDiskDistributionByInstance instanceCount=%d\n", len(instanceList))
	instanceDiskDist := make(map[int]int)
	for _, v := range instanceList {
		_, has := instanceDiskDist[v.Config.Disk]
		if !has {
			instanceDiskDist[v.Config.Disk] = 1
		} else {
			instanceDiskDist[v.Config.Disk]++
		}
	}
	instanceDiskCountList := make([]*DiskCount, 0)
	for i, v := range instanceDiskDist {
		instanceDiskCountList = append(instanceDiskCountList, &DiskCount{Disk: i, Count: v})
	}
	sort.Slice(instanceDiskCountList, func(i, j int) bool {
		return instanceDiskCountList[i].Disk > instanceDiskCountList[j].Disk
	})
	for _, v := range instanceDiskCountList {
		fmt.Printf("    disk=%5d,count=%7d\n", v.Disk, v.Count)
	}
}
