package clound

type ResourceManagement struct {
	appResourcesConfigMap    map[string]*AppResourcesConfig
	appInterferenceConfigMap map[string]map[string]int
	machineConfigMap         map[string]*MachineResourcesConfig
	machineLevelConfigPool   *MachineLevelConfigPool
	machineMap               map[string]*Machine
	machineFreePool          *MachineFreePool
	machineDeployPool        *MachineDeployPool
}

func NewResourceManagement() *ResourceManagement {
	r := &ResourceManagement{}
	r.appResourcesConfigMap = make(map[string]*AppResourcesConfig)
	r.appInterferenceConfigMap = make(map[string]map[string]int)
	r.machineConfigMap = make(map[string]*MachineResourcesConfig)
	r.machineLevelConfigPool = NewMachineLevelConfigPool()
	r.machineMap = make(map[string]*Machine)
	r.machineFreePool = NewMachineFreePool()
	r.machineDeployPool = NewMachineDeployPool()

	return r
}
