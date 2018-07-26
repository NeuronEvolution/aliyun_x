package cloud

import (
	"bytes"
	"fmt"
)

type ResourceManagement struct {
	Strategy                    Strategy
	InitialInstanceDeployConfig []*InstanceDeployConfig
	AppResourcesConfigMap       [MaxAppId]*AppResourcesConfig
	AppInterferenceConfigMap    [][MaxAppId]int
	Initializing                bool
	MachineConfigMap            [MaxMachineId]*MachineResourcesConfig
	MachineLevelConfigPool      *MachineLevelConfigPool
	MachineMap                  [MaxMachineId]*Machine
	MachineFreePool             *MachineFreePool
	MachineDeployPool           *MachineDeployPool
	DeployCommandHistory        *DeployCommandHistory
	InstanceList                [MaxInstanceId]*Instance
	InstanceDeployedMachineMap  [MaxInstanceId]*Machine

	instanceDeployedOrderByCostDescList      [MaxInstanceId]*Instance
	instanceDeployedOrderByCostDescListCount int
	instanceDeployedOrderByCostDescValid     bool
}

func NewResourceManagement() *ResourceManagement {
	r := &ResourceManagement{}
	r.Strategy = &defaultStrategy{}
	r.AppInterferenceConfigMap = make([][MaxAppId]int, MaxAppId)
	for i := 0; i < MaxAppId; i++ {
		for j := 0; j < MaxAppId; j++ {
			r.AppInterferenceConfigMap[i][j] = -1
		}
	}
	r.MachineLevelConfigPool = NewMachineLevelConfigPool()
	r.MachineFreePool = NewMachineFreePool()
	r.MachineDeployPool = NewMachineDeployPool()
	r.DeployCommandHistory = NewDeployCommandHistory()

	return r
}

func (r *ResourceManagement) DebugStatus(buf *bytes.Buffer) {
	if r.Strategy != nil {
		buf.WriteString("Strategy:")
		buf.WriteString(r.Strategy.Name())
		buf.WriteString("\n")
	}
	r.MachineDeployPool.DebugPrint(buf)

	buf.WriteString(fmt.Sprintf("cpuCost=%f,totalCommands=%d\n",
		r.CalculateTotalCostScore(), r.DeployCommandHistory.ListCount))
}

func (r *ResourceManagement) DebugPrintStatus() {
	fmt.Printf("----------------------------------------------STATUS-------------------------------------------\n")
	buf := bytes.NewBufferString("")
	r.DebugStatus(buf)
	fmt.Printf(buf.String())
	fmt.Printf("-----------------------------------------------------------------------------------------------\n")
}

func (r *ResourceManagement) SetStrategy(s Strategy) {
	r.Strategy = s
}

func (r *ResourceManagement) CalculateTotalCostScore() float64 {
	totalCost := float64(0)
	for _, m := range r.MachineDeployPool.MachineMap {
		totalCost += m.GetCostReal()
	}

	return totalCost
}
