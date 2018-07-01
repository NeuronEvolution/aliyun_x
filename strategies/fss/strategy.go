package fss

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

type FreeSmallerStrategy struct {
	R *cloud.ResourceManagement
}

func NewFreeSmallerStrategy(r *cloud.ResourceManagement) *FreeSmallerStrategy {
	s := &FreeSmallerStrategy{}
	s.R = r

	return s
}

func (s *FreeSmallerStrategy) Name() string {
	return "FreeSmallerStrategy"
}
