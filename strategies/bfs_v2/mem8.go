package bfs_v2

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
	"sort"
)

func (s *Strategy) isMem8(instance *cloud.Instance) bool {
	return instance.Config.Disk == 60 &&
		instance.Config.MemMax-instance.Config.MemAvg < 0.0000001 &&
		math.Abs(instance.Config.MemMax-8) < 0.0000001
}

func (s *Strategy) mem8(instances []*cloud.Instance) (restInstances []*cloud.Instance, err error) {
	i8s := make([]*cloud.Instance, 0)
	for _, instance := range instances {
		if s.isMem8(instance) {
			i8s = append(i8s, instance)
		}
	}

	sort.Slice(i8s, func(i, j int) bool {
		if i8s[i].Config.InferenceAppCount > i8s[j].Config.InferenceAppCount {
			return true
		} else if i8s[i].Config.InferenceAppCount == i8s[j].Config.InferenceAppCount {
			return i8s[i].Config.GetCpuDerivation() > i8s[j].Config.GetCpuDerivation()
		} else {
			return false
		}
	})

	machines := make([]*cloud.Machine, 0)
	for _, m := range s.machineDeployList {
		if m.LevelConfig.Disk == cloud.LowDisk {
			machines = append(machines, m)
		}
	}

	//deployed := make([]*cloud.Instance, 0)
	//for i, m := range machines {
	//	instance := i8s[i]
	//	if !m.ConstraintCheck(instance, 1) {
	//		continue
	//	}
	//	m.AddInstance(instance)
	//	deployed = append(deployed, instance)
	//}

	return instances, nil
}
