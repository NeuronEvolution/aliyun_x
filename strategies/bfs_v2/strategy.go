package bfs_v2

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

type BestFitStrategy struct {
	R                 *cloud.ResourceManagement
	machineDeployList []*cloud.Machine
}

func NewBestFitStrategy(r *cloud.ResourceManagement) *BestFitStrategy {
	s := &BestFitStrategy{}
	s.R = r

	return s
}

func (s *BestFitStrategy) Name() string {
	return "BestFitV2"
}
