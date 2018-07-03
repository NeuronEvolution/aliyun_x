package cloud

import (
	"math"
)

const ExpTableSize = 200000
const ExpTableStep = 10000
const ExpTableStepRev = float64(1) / float64(ExpTableStep)

var ExpTable [ExpTableSize]float64

func initExpTable() {
	for i := 0; i < len(ExpTable); i++ {
		ExpTable[i] = math.Exp(float64(i) * ExpTableStepRev)
	}
}

func Exp(r float64) float64 {
	return ExpTable[int(r*ExpTableStep)]
}
