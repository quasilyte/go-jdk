package jruntime

import (
	"sort"
)

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
	// Can use binary search because methods are sorted by name.
	i := sort.Search(len(c.Methods), func(i int) bool {
		return c.Methods[i].Name >= name
	})
	if i < len(c.Methods) && c.Methods[i].Name == name {
		return &c.Methods[i]
	}
	return nil
}
