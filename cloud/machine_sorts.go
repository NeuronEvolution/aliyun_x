package cloud

type MachineListSortByCostDesc []*Machine

func (p MachineListSortByCostDesc) Len() int {
	return len(p)
}

func (p MachineListSortByCostDesc) Less(i, j int) bool {
	return p[i].GetCost() > p[j].GetCost()
}

func (p MachineListSortByCostDesc) Swap(i, j int) {
	temp := p[i]
	p[i] = p[j]
	p[j] = temp
}
