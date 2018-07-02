package cloud

import "fmt"

type Strategy interface {
	Name() string
	PostInit() (err error)
	AddInstanceList(instanceList []*Instance) (err error)
}

type defaultStrategy struct {
}

func (s *defaultStrategy) Name() string {
	return "defaultStrategy"
}

func (s *defaultStrategy) PostInit() (err error) {
	fmt.Println("defaultStrategy.PostInit")
	return nil
}

func (s *defaultStrategy) AddInstanceList(instanceList []*Instance) (err error) {
	fmt.Println("defaultStrategy.AddInstanceList")
	return nil
}
