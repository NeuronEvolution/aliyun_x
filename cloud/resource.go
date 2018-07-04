package cloud

import "math"

type Resource struct {
	Cpu  [TimeSampleCount]float64
	Mem  [TimeSampleCount]float64
	Disk int
	P    int
	M    int
	PM   int

	ResourceCost          float64
	ResourceCostDeviation float64
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

	cpu := avgCpu / config.Cpu
	mem := avgMem / config.Mem
	disk := float64(r.Disk) / float64(config.Disk)
	p := float64(r.P) / float64(config.P)
	m := float64(r.M) / float64(config.M)
	pm := float64(r.PM) / float64(config.PM)

	r.ResourceCost = scaleCost(cpu) +
		scaleCost(mem) +
		scaleCost(disk) +
		scaleCost(p) +
		scaleCost(m) +
		scaleCost(pm)

	r.ResourceCostDeviation = calcResourceCostDeviation(cpu, mem, disk, p, m, pm)
}

func calcResourceCostDeviation(cpu float64, mem float64, disk float64, p float64, m float64, pm float64) float64 {
	avg := (cpu + mem + disk + p + m + pm) / 6
	return math.Sqrt(((cpu-avg)*(cpu-avg) + (mem-avg)*(mem-avg) + (disk-avg)*(disk-avg) +
		(p-avg)*(p-avg) + (m-avg)*(m-avg) + (pm-avg)*(pm-avg)) / 6)
}

func scaleCost(f float64) float64 {
	return Exp(4 * f)
}
