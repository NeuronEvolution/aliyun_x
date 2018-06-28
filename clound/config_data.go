package clound

type AppInterferenceConfig struct {
	AppId1       string
	AppId2       string
	Interference int
}

type AppResourcesConfig struct {
	AppId string
	Cpu   [98]float64
	Mem   [98]float64
	Disk  int
	P     int
	M     int
	PM    int

	TotalCost float64
}

func (c *AppResourcesConfig) calcTotalCost() {
	c.TotalCost = 0
}

type InstanceDeployConfig struct {
	InstanceId string
	AppId      string
	MachineId  string
}

type MachineResourcesConfig struct {
	MachineId string
	MachineLevelConfig
}
