package cloud

import "fmt"

type Strategy interface {
	AddInstance(instance *Instance) (err error)
	AddInstanceList(instanceList []*Instance) (err error)
}

type defaultStrategy struct {
}

func (s *defaultStrategy) AddInstance(instance *Instance) (err error) {
	fmt.Println("defaultStrategy.AddInstance")
	return nil
}

func (s *defaultStrategy) AddInstanceList(instanceList []*Instance) (err error) {
	fmt.Println("defaultStrategy.AddInstanceList")
	return nil
}
