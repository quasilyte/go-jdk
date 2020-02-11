package jruntime

type Class struct {
	Name    string
	Methods []Method
	Code    []byte
}

type Method struct {
	Name string
	Code []byte
}

func (c *Class) FindMethod(name string) *Method {
	for i := range c.Methods {
		if c.Methods[i].Name == name {
			return &c.Methods[i]
		}
	}
	return nil
}
