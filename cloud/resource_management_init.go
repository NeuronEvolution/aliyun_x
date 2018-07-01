package cloud

func (r *ResourceManagement) Init(
	machineResourcesConfig []*MachineResourcesConfig,
	appResourcesConfig []*AppResourcesConfig,
	appInterferenceConfig []*AppInterferenceConfig,
	instanceDeployConfig []*InstanceDeployConfig) (err error) {

	for _, v := range machineResourcesConfig {
		err = r.AddMachine(v)
		if err != nil {
			return err
		}
	}

	for _, v := range appResourcesConfig {
		err = r.SaveAppResourceConfig(v)
		if err != nil {
			return err
		}
	}

	for _, v := range appInterferenceConfig {
		err = r.SaveAppInterferenceConfig(v)
		if err != nil {
			return err
		}
	}

	err = r.InitInstanceDeploy(instanceDeployConfig)
	if err != nil {
		return err
	}

	return nil
}

func (r *ResourceManagement) ResolveAppInference() (err error) {
	return r.Strategy.ResolveAppInference()
}
