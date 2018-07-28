package bfs_v2

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
	"math"
)

func (s *Strategy) preDeploy(m *cloud.Machine, instances []*cloud.Instance) (restInstances []*cloud.Instance, err error) {
	deployed := make([]*cloud.Instance, 128)
	deployedCount := 0
	deployed[0] = instances[0]
	deployedCount++
	appCount := cloud.NewAppCountCollection()
	appCount.Add(instances[0].Config.AppId)
	resource := &cloud.Resource{}
	resource.AddResource(&instances[0].Config.Resource)

	for i := 0; i < 1; i++ {
		offset := 0
		minD := math.MaxFloat64
		var minDeployed []*cloud.Instance
		for {
			subInstances := make([]*cloud.Instance, 0)
			for ; offset < len(instances)/3; offset++ {
				if !cloud.InstancesContainsApp(subInstances, instances[offset].Config.AppId) {
					subInstances = append(subInstances, instances[offset])
					if len(subInstances) > 256 {
						break
					}
				}
			}
			if len(subInstances) == 0 {
				break
			}

			subMinD := math.MaxFloat64
			var subMinDeployed []*cloud.Instance
			for _, instance := range subInstances {
				if cloud.InstancesContainsApp(deployed[:deployedCount], instance.Config.AppId) {
					continue
				}

				if !cloud.ConstraintCheckResourceLimit(resource, &instance.Config.Resource, m.LevelConfig, 1) {
					continue
				}

				if !cloud.ConstraintCheckAppInterferenceAddInstance(instance.Config.AppId, appCount, s.R.AppInterferenceConfigMap) {
					continue
				}

				if resource.GetCostWithInstance(instance, m.LevelConfig.Cpu) > 10 {
					continue
				}

				d := resource.GetDerivationWithInstances([]*cloud.Instance{instance})
				if d >= subMinD {
					continue
				}

				subMinD = d
				subMinDeployed = []*cloud.Instance{instance}
			}

			if subMinDeployed != nil && len(subMinDeployed) > 0 {
				if subMinD < minD {
					minD = subMinD
					minDeployed = subMinDeployed

					//fmt.Println("minD", minD)
				}
			}
		}

		if minDeployed == nil || len(minDeployed) == 0 {
			break
		}

		for _, instance := range minDeployed {
			deployed[deployedCount] = instance
			deployedCount++
			appCount.Add(instance.Config.AppId)
			resource.AddResource(&instance.Config.Resource)
		}
	}

	for _, instance := range deployed[:deployedCount] {
		if !m.ConstraintCheck(instance, 1) {
			m.DebugPrint()
			panic("ConstraintCheck")
		}

		m.AddInstance(instance)
	}

	return cloud.InstancesRemove(instances, deployed[:deployedCount]), nil
}
