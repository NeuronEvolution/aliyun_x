package clound

import "fmt"

//保存应用资源需求配置
//todo 将触发调度
//todo 异步将该app相关instance拉出再重新拉入
func (r *ResourceManagement) SaveAppResourceConfig(config *AppResourcesConfig) error {
	r.appResourcesConfigMap[config.AppId] = config
	return nil
}

//保存应用冲突配置
//todo 将触发调度
//todo 异步将该app相关instance拉出再重新拉入
func (r *ResourceManagement) SaveAppInterferenceConfig(config *AppInterferenceConfig) error {
	_, hasAppResource := r.appResourcesConfigMap[config.AppId1]
	if !hasAppResource {
		return fmt.Errorf("SaveAppInterferenceConfig app %s not exists", config.AppId1)
	}

	_, hasAppResource = r.appResourcesConfigMap[config.AppId2]
	if !hasAppResource {
		return fmt.Errorf("SaveAppInterferenceConfig app %s not esists", config.AppId2)
	}

	m := r.appInterferenceConfigMap[config.AppId1]
	if m == nil {
		m = make(map[string]int, 0)
		r.appInterferenceConfigMap[config.AppId1] = m
	}
	m[config.AppId2] = config.Interference

	return nil
}
