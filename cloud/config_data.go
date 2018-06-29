package cloud

type AppInterferenceConfig struct {
	AppId1       int
	AppId2       int
	Interference int
}

type AppResourcesConfig struct {
	AppId int
	Cpu   [TimeSampleCount]float64
	Mem   [TimeSampleCount]float64
	Disk  int
	P     int
	M     int
	PM    int

	CostEval float64
}

func (c *AppResourcesConfig) calcCostEval() {
	c.CostEval = 1
}

type InstanceDeployConfig struct {
	InstanceId int
	AppId      int
	MachineId  int
}

type MachineResourcesConfig struct {
	MachineId int
	MachineLevelConfig
}
