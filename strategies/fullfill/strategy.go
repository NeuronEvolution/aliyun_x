package fullfill

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

type Strategy struct {
	R *cloud.ResourceManagement
}

func NewFullFillStrategy(r *cloud.ResourceManagement) *Strategy {
	s := &Strategy{}
	s.R = r

	return s
}

func (s *Strategy) Name() string {
	return "FullFillStrategy"
}
