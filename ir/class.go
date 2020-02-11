package ir

type Class struct {
	Name    string
	PkgName string
	Methods []Method
}

type Method struct {
	Name string
	Code []Inst
}
