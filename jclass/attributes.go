package jclass

type Attribute interface {
	attribute()
}

type (
	RawAttribute struct {
		NameIndex uint16
		Data      []byte
	}

	CodeAttribute struct {
		MaxStack       uint16
		MaxLocals      uint16
		Code           []byte
		ExceptionTable []ExceptionHandler
		Attrs          []Attribute
	}

	StackMapTableAttribute struct {
		Frames []StackMapFrame
	}
)

func (RawAttribute) attribute()           {}
func (CodeAttribute) attribute()          {}
func (StackMapTableAttribute) attribute() {}
