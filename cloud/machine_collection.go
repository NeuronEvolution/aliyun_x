package cloud

import "fmt"

type MachineCollection struct {
	Map       map[int]*Machine
	List      []*Machine
	ListCount int
}

func NewMachineCollection() *MachineCollection {
	c := &MachineCollection{}
	c.Map = make(map[int]*Machine)
	c.List = make([]*Machine, 0)

	return c
}

func (c *MachineCollection) Add(m *Machine) {
	//debugLog("MachineCollection.Add %s", m.MachineId)
	_, has := c.Map[m.MachineId]
	if has {
		panic(fmt.Errorf("MachineCollection.Add %s exists", m.MachineId))
	}

	c.Map[m.MachineId] = m
	c.List = append(c.List, m)
	c.ListCount++
}

func (c *MachineCollection) Remove(machineId int) {
	debugLog("MachineCollection.Remove %s", machineId)
	_, has := c.Map[machineId]
	if !has {
		panic(fmt.Errorf("MachineCollection.Add %s not exists", machineId))
	}

	delete(c.Map, machineId)
	for i := 0; i < c.ListCount; i++ {
		v := c.List[i]
		if v.MachineId == machineId {
			c.List[i] = nil
			if c.ListCount > 1 && i < c.ListCount-1 {
				c.List[i] = c.List[c.ListCount-1]
			}
			c.ListCount--
		}
	}
}

func (c *MachineCollection) Pop() (m *Machine) {
	//debugLog("MachineCollection.Pop ListCount=%d", c.ListCount)
	if c.ListCount == 0 {
		return nil
	}

	c.ListCount--

	m = c.List[c.ListCount]
	c.List[c.ListCount] = nil
	delete(c.Map, m.MachineId)

	if c.ListCount != len(c.Map) {
		panic("111")
	}

	//debugLog("MachineCollection.Pop success %s", m.MachineId)

	return m
}

func (c *MachineCollection) Has(machinedId int) bool {
	_, has := c.Map[machinedId]
	return has
}

func (c *MachineCollection) First() *Machine {
	if c.ListCount == 0 {
		return nil
	}

	return c.List[0]
}

func (c *MachineCollection) Last() *Machine {
	if c.ListCount == 0 {
		return nil
	}

	return c.List[c.ListCount-1]
}
