package cloud

type ResourceManagement struct {
	Strategy                 Strategy
	AppResourcesConfigMap    map[int]*AppResourcesConfig
	AppInterferenceConfigMap [][MaxAppId]int
	MachineConfigMap         map[int]*MachineResourcesConfig
	MachineLevelConfigPool   *MachineLevelConfigPool
	MachineMap               map[int]*Machine
	MachineFreePool          *MachineFreePool
	MachineDeployPool        *MachineDeployPool
}

func NewResourceManagement() *ResourceManagement {
	r := &ResourceManagement{}
	r.Strategy = &defaultStrategy{}
	r.AppResourcesConfigMap = make(map[int]*AppResourcesConfig)
	r.AppInterferenceConfigMap = make([][MaxAppId]int, MaxAppId)
	for i := 0; i < MaxAppId; i++ {
		for j := 0; j < MaxAppId; j++ {
			r.AppInterferenceConfigMap[i][j] = -1
		}
	}
	r.MachineConfigMap = make(map[int]*MachineResourcesConfig)
	r.MachineLevelConfigPool = NewMachineLevelConfigPool()
	r.MachineMap = make(map[int]*Machine)
	r.MachineFreePool = NewMachineFreePool()
	r.MachineDeployPool = NewMachineDeployPool()

	return r
}

func (r *ResourceManagement) DebugDeployStatus() {
	r.MachineDeployPool.DebugPrint()
}

func (r *ResourceManagement) SetStrategy(s Strategy) {
	r.Strategy = s
}

func (r *ResourceManagement) CalculateTotalCostScore() float64 {
	totalCost := float64(0)
	for _, m := range r.MachineDeployPool.MachineMap {
		totalCost += m.CalculateCost()
	}

	return totalCost
}
