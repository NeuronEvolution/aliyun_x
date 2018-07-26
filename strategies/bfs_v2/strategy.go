package bfs_v2

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

type Strategy struct {
	R                 *cloud.ResourceManagement
	machineDeployList []*cloud.Machine
}

func NewStrategy(r *cloud.ResourceManagement) *Strategy {
	s := &Strategy{}
	s.R = r

	return s
}

func (s *Strategy) Name() string {
	return "BestFitV2"
}
