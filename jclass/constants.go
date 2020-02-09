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
		Value int32
	}

	LongConst struct {
		Value int64
	}

	FloatConst struct {
		Value float32
	}

	DoubleConst struct {
		Value float64
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
func (*LongConst) constant()        {}
func (*FloatConst) constant()       {}
func (*DoubleConst) constant()      {}
func (*ClassConst) constant()       {}
func (*FieldrefConst) constant()    {}
func (*MethodrefConst) constant()   {}
func (*NameAndTypeConst) constant() {}
