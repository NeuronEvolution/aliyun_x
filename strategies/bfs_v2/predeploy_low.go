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

func (s *Strategy) preDeployLow(instances []*cloud.Instance, inferenceLimit int) (restInstances []*cloud.Instance, err error) {
	i8s := make([]*cloud.Instance, 0)
	for _, instance := range instances {
		if s.isMem8(instance) {
			i8s = append(i8s, instance)
		}
	}

	sort.Slice(i8s, func(i, j int) bool {
		c1 := i8s[i].Config
		c2 := i8s[j].Config
		if c1.InferenceAppCount < inferenceLimit && c2.InferenceAppCount < inferenceLimit {
			return c1.GetCpuDerivation() > c2.GetCpuDerivation()
		}

		if c1.InferenceAppCount > c2.InferenceAppCount {
			return true
		} else if c1.InferenceAppCount == c2.InferenceAppCount {
			return c1.GetCpuDerivation() > c2.GetCpuDerivation()
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

	i8sRest := i8s
	for _, m := range machines {
		i8sRest, err = s.preDeployLowMachine(m, i8sRest)
		if err != nil {
			return nil, err
		}
	}

	i8sDeployed := cloud.InstancesRemove(i8s, i8sRest)

	return cloud.InstancesRemove(instances, i8sDeployed), nil
}

func (s *Strategy) preDeployLowMachine(m *cloud.Machine, instances []*cloud.Instance) (restInstances []*cloud.Instance, err error) {
	for i := 0; i < 8; i++ {
		for _, instance := range instances {
			if instance.Config.CpuAvg > 6 || instance.Config.CpuMax > 8 {
				continue
			}

			if cloud.InstancesContainsApp(m.InstanceArray[:m.InstanceArrayCount], instance.Config.AppId) {
				continue
			}

		}
	}

	return instances, nil
}
