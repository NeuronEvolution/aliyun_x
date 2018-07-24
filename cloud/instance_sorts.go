package cloud

import "sort"

//92,288,1024,7,7,9
func SortInstanceByTotalMax(p []*Instance) {
	sort.Slice(p, func(i, j int) bool {
		c1 := p[i].Config
		c2 := p[j].Config
		a1 := float64(c1.Disk)/float64(1024) + c1.CpuMax/float64(46) + c1.MemMax/float64(288) +
			float64(c1.P)/float64(7) + float64(c1.M)/float64(7) + float64(c1.PM)/float64(9)
		a2 := float64(c2.Disk)/float64(1024) + c2.CpuMax/float64(46) + c2.MemMax/float64(288) +
			float64(c2.P)/float64(7) + float64(c2.M)/float64(7) + float64(c2.PM)/float64(9)
		return a1 > a2
	})
}

func SortInstanceByDisk(p []*Instance) {
	sort.Slice(p, func(i, j int) bool {
		if p[i].Config.Disk > p[j].Config.Disk {
			return true
		} else if p[i].Config.Disk == p[j].Config.Disk {
			a1 := p[i].Config.CpuMax*float64(288) + p[i].Config.MemMax*float64(46)
			a2 := p[j].Config.CpuMax*float64(288) + p[j].Config.MemMax*float64(46)
			if a1 > a2 {
				return true
			} else if a1 == a2 {
				if p[i].Config.AppId < p[j].Config.AppId {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	})
}

func SortInstanceByCpu(p []*Instance) {
	sort.Slice(p, func(i, j int) bool {
		if p[i].Config.CpuMax > p[j].Config.CpuMax {
			return true
		} else if p[i].Config.CpuMax > p[j].Config.CpuMax {
			a1 := p[i].Config.MemMax*float64(1024) + float64(p[i].Config.Disk)*float64(288)
			a2 := p[j].Config.MemMax*float64(1024) + float64(p[j].Config.Disk)*float64(288)
			if a1 > a2 {
				return true
			} else if a1 == a2 {
				if p[i].Config.AppId < p[j].Config.AppId {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	})
}

func SortInstanceByMem(p []*Instance) {
	sort.Slice(p, func(i, j int) bool {
		if p[i].Config.MemMax > p[j].Config.MemMax {
			return true
		} else if p[i].Config.MemMax > p[j].Config.MemMax {
			a1 := p[i].Config.CpuMax*float64(1024) + float64(p[i].Config.Disk)*float64(46)
			a2 := p[j].Config.CpuMax*float64(1024) + float64(p[j].Config.Disk)*float64(46)
			if a1 > a2 {
				return true
			} else if a1 == a2 {
				if p[i].Config.AppId < p[j].Config.AppId {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	})
}
