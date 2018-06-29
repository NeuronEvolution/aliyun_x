package cloud

import "fmt"

type AppCount struct {
	AppId int
	Count int
}

type AppCountCollection struct {
	List      [64]AppCount
	ListCount int
}

func NewAppCountCollection() *AppCountCollection {
	c := &AppCountCollection{}

	return c
}

func (c *AppCountCollection) Add(appId int) {
	for i := 0; i < c.ListCount; i++ {
		if c.List[i].AppId == appId {
			c.List[i].Count++
			return
		}
	}

	item := &c.List[c.ListCount]
	item.AppId = appId
	item.Count = 1
	c.ListCount++
}

func (c *AppCountCollection) Remove(appId int) {
	for i := 0; i < c.ListCount; i++ {
		item := &c.List[i]
		if item.AppId == appId {
			if item.Count <= 0 {
				panic(fmt.Errorf("AppCountCollection.Remove appId %d count<=0", appId))
			}

			item.Count--
			if item.Count == 0 {
				if i != c.ListCount-1 {
					last := &c.List[c.ListCount-1]
					item.AppId = last.AppId
					item.Count = last.Count
				}
			}

			c.ListCount--
		}
	}

	panic(fmt.Errorf("AppCountCollection.Remove appId %d not exists", appId))
}
