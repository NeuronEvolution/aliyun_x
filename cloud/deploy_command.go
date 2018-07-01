package cloud

type DeployCommand struct {
	AppId      int
	InstanceId int
	MachineId  int
}

type DeployCommandHistory struct {
	List      []DeployCommand
	ListCount int
}

func NewDeployCommandHistory() *DeployCommandHistory {
	h := &DeployCommandHistory{}
	h.List = make([]DeployCommand, MaxDeployCommandCount)

	return h
}

func (h *DeployCommandHistory) Push(appId int, instanceId int, machineId int) {
	//debugLog("DeployCommandHistory.Push %d %d",InstanceId,machineId)

	item := &h.List[h.ListCount]
	item.AppId = appId
	item.InstanceId = instanceId
	item.MachineId = machineId
	h.ListCount++
}

func (h *DeployCommandHistory) DebugPrint() {
	debugLog("DeployCommandHistory.DebugPrint")
	for _, v := range h.List[:h.ListCount] {
		debugLog("%d.%d -> %d", v.AppId, v.InstanceId, v.MachineId)
	}
	debugLog("DeployCommandHistory.DebugPrint totalCount=%d", h.ListCount)
}
