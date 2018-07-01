package cloud

import (
	"bytes"
	"fmt"
)

type ResourceManagement struct {
	Strategy                    Strategy
	InitialInstanceDeployConfig []*InstanceDeployConfig
	AppResourcesConfigMap       map[int]*AppResourcesConfig
	AppInterferenceConfigMap    [][MaxAppId]int
	MachineConfigMap            map[int]*MachineResourcesConfig
	MachineLevelConfigPool      *MachineLevelConfigPool
	MachineMap                  map[int]*Machine
	MachineFreePool             *MachineFreePool
	MachineDeployPool           *MachineDeployPool
	DeployCommandHistory        *DeployCommandHistory
	InstanceList                [MaxInstanceId]*Instance
	InstanceMachineMap          [MaxInstanceId]*Machine
}

func NewResourceManagement() *ResourceManagement {
	r := &ResourceManagement{}
	r.Strategy = &defaultStrategy{}
	r.AppResourcesConfigMap = make(map[int]*AppResourcesConfig)
	r.AppInterferenceConfigMap = make([][MaxAppId]int, MaxAppId)
	for i := 0; i < MaxAppId; i++ {
		for j := 0; j < MaxAppId; j++ {
			r.AppInterferenceConfigMap[i][j] = -1
		}
	}
	r.MachineConfigMap = make(map[int]*MachineResourcesConfig)
	r.MachineLevelConfigPool = NewMachineLevelConfigPool()
	r.MachineMap = make(map[int]*Machine)
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

	buf.WriteString(fmt.Sprintf("cost=%f,totalCommands=%d\n",
		r.CalculateTotalCostScore(), r.DeployCommandHistory.ListCount))
}

func (r *ResourceManagement) DebugPrintStatus() {
	fmt.Printf("-----------------------------------------------------------------\n")
	buf := bytes.NewBufferString("")
	r.DebugStatus(buf)
	fmt.Printf(buf.String())
	fmt.Printf("#################################################################\n")
}

func (r *ResourceManagement) SetStrategy(s Strategy) {
	r.Strategy = s
}

func (r *ResourceManagement) CalculateTotalCostScore() float64 {
	totalCost := float64(0)
	for _, m := range r.MachineDeployPool.MachineMap {
		totalCost += m.CalculateCost()
	}

	return totalCost
}
