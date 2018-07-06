package cloud

func constraintCheckAppInterference(c *AppCountCollection, m [][MaxAppId]int) bool {
	//debugLog("constraintCheckAppInterference %v", c.List[:c.ListCount])
	for _, v1 := range c.List[:c.ListCount] {
		for _, v2 := range c.List[:c.ListCount] {
			maxCount := m[v1.AppId][v2.AppId]
			if maxCount != -1 {
				if v1.AppId == v2.AppId {
					maxCount++
				}
				if v2.Count > maxCount {
					return false
				}
			}
		}
	}

	return true
}

//app冲突约束检测
//appId:要新增的appId
//c:当前机器已部署的每个app的数量
//m:冲突配置
func constraintCheckAppInterferenceAddInstance(appId int, c *AppCountCollection, m [][MaxAppId]int) bool {
	//debugLog("constraintCheckAppInterferenceAddInstance appId=%d %v", appId, c.List[:c.ListCount])
	appCount := 0
	for _, v := range c.List[:c.ListCount] {
		if v.AppId == appId {
			appCount = v.Count
			break
		}
	}
	appCount++

	//<appId,AppId>
	maxCount := m[appId][appId]
	if maxCount != -1 && appCount > maxCount+1 {
		debugLog("constraintCheckAppInterferenceAddInstance 1 failed app=%d,count2=%d,max=%d",
			appId, appCount, maxCount)
		return false
	}

	for _, v := range c.List[:c.ListCount] {
		if v.AppId == appId {
			continue
		}

		//<appIdOther,appId>
		maxCount := m[v.AppId][appId]
		if maxCount != -1 && appCount > maxCount {
			debugLog("constraintCheckAppInterferenceAddInstance 2 failed app1=%d,app2=%d,count2=%d,max=%d",
				v.AppId, appId, appCount, maxCount)
			return false
		}

		//<appId,appIdOther>
		if appCount == 1 { //已经存在的app，数量增加不影响冲突结果
			maxCount = m[appId][v.AppId]
			if maxCount != -1 && v.Count > maxCount {
				debugLog("constraintCheckAppInterferenceAddInstance 3 failed app1=%d,app2=%d,count2=%d,max=%d",
					appId, v.AppId, v.Count, maxCount)
				return false
			}
		}
	}

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
		debugLog("constraintCheckResourceLimit failed P %d %d %d", m.P, i.P, c.P)
		return false
	}

	if m.M+i.M > c.M {
		debugLog("constraintCheckResourceLimit failed M %d %d %d", m.M, i.M, c.M)
		return false
	}

	if m.PM+i.PM > c.PM {
		debugLog("constraintCheckResourceLimit failed PM %d %d %d", m.PM, i.PM, c.PM)
		return false
	}

	for index, v := range m.Cpu {
		if v+i.Cpu[index] > c.Cpu {
			debugLog("constraintCheckResourceLimit failed Cpu %d %f %f %f", index, v, i.Cpu[index], c.Cpu)
			return false
		}
	}

	for index, v := range m.Mem {
		if v+i.Mem[index] > c.Mem {
			debugLog("constraintCheckResourceLimit failed Mem %d %f %f %f", index, v, i.Mem[index], c.Mem)
			return false
		}
	}

	return true
}
