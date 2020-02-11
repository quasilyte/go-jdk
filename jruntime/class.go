package jruntime

type Class struct {
	Name string
}

type Method struct{}

func (c *Class) FindMethod(name string) *Method {
	// This is a stub.
	return &Method{}
}
