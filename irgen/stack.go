package irgen

// operandStack simulates JVM run-time operands stack.
type operandStack struct {
	values []stackValue
}

// reset clears a stack.
// Memory is re-used.
func (st *operandStack) reset() {
	st.values = st.values[:0]
}

func (st *operandStack) push(kind valueKind, v int64) {
	st.values = append(st.values, stackValue{kind: kind, value: v})
}

// top returns last pushed stack value.
func (st *operandStack) top() stackValue {
	return st.values[len(st.values)-1]
}

// drop removes n top values from a stack.
func (st *operandStack) drop(n int) {
	st.values = st.values[:len(st.values)-n]
}

// get returns n-th stack value.
// Indexing starts from the top, get(0) is identical to top().
func (st *operandStack) get(n int) stackValue {
	return st.values[len(st.values)-n-1]
}

type stackValue struct {
	kind  valueKind
	value int64
}

type valueKind int

const (
	valueInvalid valueKind = iota

	valueIntConst
	valueLongConst
	valueFloatConst
	valueDoubleConst

	valueIntLocal
	valueLongLocal
	valueFloatLocal
	valueDoubleLocal

	valueTmp
	valueFlags
)
