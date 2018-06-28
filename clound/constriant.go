package clound

func constraintCheckAppInterference(appId string, appCountMap map[string]int, interferenceMap map[string]map[string]int) bool {
	debugLog("constraintCheckAppInterference %s %v", appId, appCountMap)
	appCount := appCountMap[appId]
	appCount++

	//<appIdOther,appId>
	for appIdOther := range appCountMap {
		if appIdOther == appId {
			continue
		}

		m := interferenceMap[appIdOther]
		if m == nil {
			continue
		}

		maxCount, has := m[appId]
		if !has {
			continue
		}

		if appCount > maxCount {
			debugLog("constraintCheckAppInterference failed app1=%s,app2=%s,count2=%d,max=%d",
				appIdOther, appId, appCount, maxCount)
			return false
		}
	}

	m := interferenceMap[appId]
	if m != nil {
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
				debugLog("constraintCheckAppInterference failed app1=%s,app2=%s,count2=%d,max=%d",
					appId, appIdOther, countOther, maxCount)
				return false
			}
		}

		//<appId,appId>
		maxCount, has := m[appId]
		if has {
			if appCount > maxCount {
				debugLog("constraintCheckAppInterference failed app1=%s,app2=%s,count2=%d,max=%d",
					appId, appId, appCount, maxCount)
				return false
			}
		}
	}

	return true
}

func constraintCheckResourceLimit(m *Machine, instanceList []*Instance) bool {

	return true
}
