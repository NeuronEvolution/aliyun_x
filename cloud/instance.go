package cloud

type Instance struct {
	ResourceManagement *ResourceManagement
	InstanceId         int
	Config             *AppResourcesConfig
	ResourceCost       float64
}

func NewInstance(r *ResourceManagement, instanceId int, config *AppResourcesConfig) *Instance {
	i := &Instance{}
	i.ResourceManagement = r
	i.InstanceId = instanceId
	i.Config = config
	i.ResourceCost = config.ResourceCost

	return i
}

func CopyInstanceList(p []*Instance) (r []*Instance) {
	if p == nil {
		return nil
	}

	r = make([]*Instance, len(p))
	for i, v := range p {
		r[i] = v
	}

	return r
}
