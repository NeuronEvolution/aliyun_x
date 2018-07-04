package cloud

func (r *ResourceManagement) Init(
	machineResourcesConfig []*MachineResourcesConfig,
	appResourcesConfig []*AppResourcesConfig,
	appInterferenceConfig []*AppInterferenceConfig,
	instanceDeployConfig []*InstanceDeployConfig) (err error) {

	r.Initializing = true
	defer func() { r.Initializing = false }()

	if machineResourcesConfig != nil {
		for _, v := range machineResourcesConfig {
			err = r.AddMachine(v)
			if err != nil {
				return err
			}
		}
	}

	if appResourcesConfig != nil {
		for _, v := range appResourcesConfig {
			err = r.SaveAppResourceConfig(v)
			if err != nil {
				return err
			}
		}
	}

	if appInterferenceConfig != nil {
		for _, v := range appInterferenceConfig {
			err = r.SaveAppInterferenceConfig(v)
			if err != nil {
				return err
			}
		}
	}

	if instanceDeployConfig != nil {
		err = r.InitInstanceDeploy(instanceDeployConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ResourceManagement) PostInit() (err error) {
	return r.Strategy.PostInit()
}
