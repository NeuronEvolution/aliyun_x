package clound

type Instance struct {
	ResourceManagement *ResourceManagement
	InstanceId         string
	Config             *AppResourcesConfig
}

func NewInstance(r *ResourceManagement, instanceId string, config *AppResourcesConfig) *Instance {
	i := &Instance{}
	i.ResourceManagement = r
	i.InstanceId = instanceId
	i.Config = config

	return i
}

type InstanceArray []*Instance

func (p InstanceArray) Len() int {
	return len(p)
}

func (p InstanceArray) Less(i, j int) bool {
	return p[i].Config.TotalCost > p[j].Config.TotalCost
}

func (p InstanceArray) Swap(i, j int) {
	temp := p[i]
	p[i] = p[j]
	p[j] = temp
}
