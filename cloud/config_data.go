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
	avgCpu := float64(0)
	for _, v := range c.Cpu {
		avgCpu += v
	}
	avgCpu = avgCpu / float64(len(c.Cpu))

	avgMem := float64(0)
	for _, v := range c.Mem {
		avgMem += v
	}
	avgMem = avgMem / float64(len(c.Mem))

	c.CostEval = avgCpu/MachineCpuMax +
		avgMem/MachineMemMax +
		float64(c.Disk)/MachineDiskMax +
		float64(c.P)/MachinePMax +
		float64(c.M)/MachineMMax +
		float64(c.PM)/MachinePMMax
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
