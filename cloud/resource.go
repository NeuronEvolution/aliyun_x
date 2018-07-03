package cloud

type Resource struct {
	Cpu  [TimeSampleCount]float64
	Mem  [TimeSampleCount]float64
	Disk int
	P    int
	M    int
	PM   int

	ResourceCost float64
}

func (r *Resource) calcCostEval(config *MachineLevelConfig) {
	avgCpu := float64(0)
	for _, v := range r.Cpu {
		avgCpu += v
	}
	avgCpu = avgCpu / float64(len(r.Cpu))

	avgMem := float64(0)
	for _, v := range r.Mem {
		avgMem += v
	}
	avgMem = avgMem / float64(len(r.Mem))

	r.ResourceCost = avgCpu/config.Cpu +
		avgMem/config.Mem +
		float64(r.Disk)/float64(config.Disk) +
		float64(r.P)/float64(config.P) +
		float64(r.M)/float64(config.M) +
		float64(r.PM)/float64(config.PM)
}
