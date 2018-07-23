package fullfill

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

func (s *Strategy) AddInstanceList(instanceList []*cloud.Instance) (err error) {
	instanceOrderbyDisk := make([]*cloud.Instance, len(instanceList))
	instanceOrderbyCpu := make([]*cloud.Instance, len(instanceList))
	instanceOrderbyMem := make([]*cloud.Instance, len(instanceList))
	instanceOrderbyP := make([]*cloud.Instance, len(instanceList))
	instanceOrderbyM := make([]*cloud.Instance, len(instanceList))
	instanceOrderbyPM := make([]*cloud.Instance, len(instanceList))

	return nil
}
