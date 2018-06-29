package cloud

type Instance struct {
	ResourceManagement *ResourceManagement
	InstanceId         int
	Config             *AppResourcesConfig
	CostEval           float64
}

func NewInstance(r *ResourceManagement, instanceId int, config *AppResourcesConfig) *Instance {
	i := &Instance{}
	i.ResourceManagement = r
	i.InstanceId = instanceId
	i.Config = config
	i.CostEval = config.CostEval

	return i
}

type InstanceArray []*Instance

func (p InstanceArray) Len() int {
	return len(p)
}

func (p InstanceArray) Less(i, j int) bool {
	return p[i].CostEval > p[j].CostEval
}

func (p InstanceArray) Swap(i, j int) {
	temp := p[i]
	p[i] = p[j]
	p[j] = temp
}
