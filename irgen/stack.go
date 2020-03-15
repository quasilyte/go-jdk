package irgen

// operandStack simulates JVM run-time operands stack.
type operandStack struct {
	tmp      int64
	values   []stackValue
	freelist []int64
}

// reset clears a stack.
// Memory is re-used.
func (st *operandStack) reset() {
	st.tmp = 0
	st.values = st.values[:0]
	st.freelist = st.freelist[:0]
}

func (st *operandStack) nextTmp() int64 {
	if len(st.freelist) != 0 {
		i := st.freelist[len(st.freelist)-1]
		st.freelist = st.freelist[:len(st.freelist)-1]
		return i
	}
	v := st.tmp
	st.tmp++
	return v
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
	for i := 0; i < n; i++ {
		v := st.get(i)
		if v.kind == valueTmp {
			st.freelist = append(st.freelist, v.value)
		}
	}

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
