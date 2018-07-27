package bfs_v2

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
)

func (s *Strategy) preDeploy(m *cloud.Machine, instances []*cloud.Instance) (restInstances []*cloud.Instance, err error) {
	deployed := make([]*cloud.Instance, 0)

	m.AddInstance(instances[0])
	deployed = append(deployed, instances[0])

	minD := math.MaxFloat64
	var minInstance *cloud.Instance
	for _, instance := range instances[1 : len(instances)/3] {
		d := m.GetDerivationWithInstance(instance)
		if d >= minD {
			continue
		}

		if !m.ConstraintCheck(instance, 1) {
			continue
		}

		minD = d
		minInstance = instance
		//fmt.Println("minD", minD)

	}

	if minInstance != nil {
		m.AddInstance(minInstance)
		deployed = append(deployed, minInstance)
	}

	return cloud.InstancesRemove(instances, deployed), nil
}
