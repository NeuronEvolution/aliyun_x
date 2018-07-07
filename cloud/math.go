package cloud

import (
	"math"
)

const ExpTableSize = 200000
const ExpTableStep = 10000
const ExpTableStepRev = float64(1) / float64(ExpTableStep)

const SqrtTableSize = 200000
const SqrtTableStep = 10000
const SqrtTableStepRev = float64(1) / float64(SqrtTableStep)

var ExpTable [ExpTableSize]float64
var SqrtTable [SqrtTableSize]float64

func initExpTable() {
	for i := 0; i < len(ExpTable); i++ {
		ExpTable[i] = math.Exp(float64(i) * ExpTableStepRev)
	}
}

func initSqrtTable() {
	for i := 0; i < len(SqrtTable); i++ {
		SqrtTable[i] = math.Sqrt(float64(i) * SqrtTableStepRev)
	}
}

func Exp(r float64) float64 {
	return ExpTable[int(r*ExpTableStep)]
}

func Sqrt(r float64) float64 {
	return r
}

func Statistics(arr [TimeSampleCount]float64) (avg float64, dev float64, min float64, max float64) {
	min = math.MaxFloat64
	max = -math.MaxFloat64
	for _, v := range arr {
		avg += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	avg = avg / float64(len(arr))

	for _, v := range arr {
		dev += (v - avg) * (v - avg)
	}
	dev = math.Sqrt(dev / float64(len(arr)))

	return
}
