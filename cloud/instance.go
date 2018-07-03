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

type InstanceListSortByCostEvalDesc []*Instance

func (p InstanceListSortByCostEvalDesc) Len() int {
	return len(p)
}

func (p InstanceListSortByCostEvalDesc) Less(i, j int) bool {
	return p[i].ResourceCost > p[j].ResourceCost
}

func (p InstanceListSortByCostEvalDesc) Swap(i, j int) {
	temp := p[i]
	p[i] = p[j]
	p[j] = temp
}
