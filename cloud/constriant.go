package cloud

import "fmt"

//app冲突约束检测
//appId:要新增的appId
//appCountMap:当前机器已部署的每个app的数量
//interferenceMap:冲突配置
func constraintCheckAppInterference(appId int, appCountMap map[int]int, interferenceMap []map[int]int) bool {
	//debugLog("constraintCheckAppInterference %s %v", appId, appCountMap)
	appCount := appCountMap[appId]
	appCount++

	if len(appCountMap) > 20 {
		fmt.Println("appCountMap", len(appCountMap))
	}

	//<appIdOther,appId>
	for appIdOther := range appCountMap {
		if appIdOther == appId {
			continue
		}

		/*m := interferenceMap[appIdOther]
		if m == nil {
			continue
		}

		maxCount, has := m[appId]
		if !has {
			continue
		}

		if appCount > maxCount {
			//debugLog("constraintCheckAppInterference failed app1=%s,app2=%s,count2=%d,max=%d",
			//	appIdOther, appId, appCount, maxCount)
			return false
		}*/
	}

	//m := interferenceMap[appId]
	//if m != nil {
	/*
		//<appId,appIdOther>
		for appIdOther, countOther := range appCountMap {
			if appIdOther == appId {
				continue
			}

			maxCount, has := m[appIdOther]
			if !has {
				continue
			}

			if countOther > maxCount {
				//debugLog("constraintCheckAppInterference failed app1=%s,app2=%s,count2=%d,max=%d",
				//appId, appIdOther, countOther, maxCount)
				return false
			}
		}

		//<appId,appId>
		maxCount, has := m[appId]
		if has {
			if appCount > maxCount+1 {
				//debugLog("constraintCheckAppInterference failed app1=%s,app2=%s,count2=%d,max=%d",
				//appId, appId, appCount, maxCount)
				return false
			}
		}*/
	//}

	return true
}

func constraintCheckResourceLimit(m *Machine, instance *Instance) bool {
	c := m.LevelConfig
	i := instance.Config

	if m.Disk+i.Disk > c.Disk {
		//debugLog("constraintCheckResourceLimit failed Disk %d %d %d", m.Disk, i.Disk, c.Disk)
		return false
	}

	if m.P+i.P > c.P {
		//debugLog("constraintCheckResourceLimit failed P %d %d %d", m.P, i.P, c.P)
		return false
	}

	if m.M+i.M > c.M {
		//debugLog("constraintCheckResourceLimit failed M %d %d %d", m.M, i.M, c.M)
		return false
	}

	if m.PM+i.PM > c.PM {
		//debugLog("constraintCheckResourceLimit failed PM %d %d %d", m.PM, i.PM, c.PM)
		return false
	}

	for index, v := range m.Cpu {
		if v+i.Cpu[index] > c.Cpu {
			//debugLog("constraintCheckResourceLimit failed Cpu %d %d %d %d", index, v, i.Cpu[index], c.Cpu)
			return false
		}
	}

	for index, v := range m.Mem {
		if v+i.Mem[index] > c.Mem {
			//debugLog("constraintCheckResourceLimit failed Mem %d %d %d %d", index, v, i.Mem[index], c.Mem)
			return false
		}
	}

	return true
}
