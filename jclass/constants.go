package jclass

type Const interface {
	constant()
}

type (
	Utf8Const struct {
		Value string
	}

	ClassConst struct {
		Name string
	}

	IntConst struct {
		Value int64
	}

	FieldrefConst struct {
		ClassName  string
		Name       string
		Descriptor string
	}

	MethodrefConst struct {
		ClassName  string
		Name       string
		Descriptor string
	}

	NameAndTypeConst struct {
		Name       string
		Descriptor string
	}
)

func (*Utf8Const) constant()        {}
func (*IntConst) constant()         {}
func (*ClassConst) constant()       {}
func (*FieldrefConst) constant()    {}
func (*MethodrefConst) constant()   {}
func (*NameAndTypeConst) constant() {}
