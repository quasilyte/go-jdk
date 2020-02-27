package irgen

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/quasilyte/go-jdk/ir"
	"github.com/quasilyte/go-jdk/irfmt"
	"github.com/quasilyte/go-jdk/loader"
	"github.com/quasilyte/go-jdk/vmdat"
)

func TestArgCount(t *testing.T) {
	tests := []struct {
		d    string
		want int
	}{
		{"()V", 0},
		{"()I", 0},
		{"(I)I", 1},
		{"(II)I", 2},
		{"(LFoo;)I", 1},
		{"([[I)I", 1},
		{"([I[I)I", 2},
		{"(LFoo;LBar;)I", 2},
		{"(LObject;I)I", 2},
		{"(Z[LObject;Z)V", 3},
	}

	for _, test := range tests {
		have := argsCount(test.d)
		if have != test.want {
			t.Errorf("argsCount(%q):\nhave: %d\nwant: %d", test.d, have, test.want)
		}
	}
}

func TestIrgen(t *testing.T) {
	absTestdata, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatalf("abs(testdata): %v", err)
	}

	// Compile Java file.
	{
		args := []string{
			"-cp", "testdata",
			"testdata/irgen/C1.java",
		}
		out, err := exec.Command("javac", args...).CombinedOutput()
		if err != nil {
			t.Fatalf("javac: %v: %s", err, out)
		}
	}

	// Collect magic comments.
	type irTest struct {
		lineNum int
		irText  string
	}
	tests := map[string]irTest{}
	{
		methodNameRE := regexp.MustCompile(`(\w+)\(`)
		data, err := ioutil.ReadFile(filepath.Join("testdata", "irgen", "C1.java"))
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(string(data), "\n")
		i := 0
		for i < len(lines) {
			l := lines[i]
			if strings.Contains(l, "// slots=") {
				var irLines []string
				for {
					comment := strings.Index(lines[i], "// ")
					if comment == -1 {
						break
					}
					irline := lines[i][comment+len("// "):]
					irLines = append(irLines, irline)
					i++
				}
				methodName := methodNameRE.FindStringSubmatch(lines[i])[1]
				tests[methodName] = irTest{
					lineNum: i + 1,
					irText:  strings.Join(irLines, "\n") + "\n",
				}
			} else {
				i++
			}
		}
	}

	var st vmdat.State
	st.Init()
	packages, err := loader.LoadPackage(&st, "irgen", &loader.Config{
		ClassPath: []string{absTestdata},
	})
	if err != nil {
		t.Fatalf("load package: %v", err)
	}
	if len(packages) != 1 {
		t.Fatalf("expected 1 package, got %d packages", len(packages))
	}
	if err := Generate(&st, packages); err != nil {
		t.Fatalf("irgen: %v", err)
	}

	pkg := packages[0]
	for i := range pkg.Classes {
		c := &pkg.Classes[i]
		for j := range c.Methods {
			m := &c.Methods[j]
			test, ok := tests[m.Out.Name]
			if !ok {
				continue
			}
			have := sprintMethod(&st, m)
			want := test.irText
			if have != want {
				t.Errorf("C1.java:%d: method %s:\nhave:\n%s\nwant:\n%s",
					test.lineNum, m.Out.Name, have, want)
			}
		}
	}
}

var branchArgRE = regexp.MustCompile(`@(\d+)`)

func sprintMethod(st *vmdat.State, m *ir.Method) string {
	// TODO: move to irfmt package and use in javap as well.

	var buf strings.Builder

	index2label := map[int]string{}
	for i, inst := range m.Code {
		if inst.Flags.IsJumpTarget() && i != 0 {
			index2label[i] = fmt.Sprintf("label%d", len(index2label))
		}
	}

	fmt.Fprintf(&buf, "slots=%d\n", m.Out.FrameSlots)
	blockIndex := -1
	for i, inst := range m.Code {
		if inst.Flags.IsBlockLead() {
			blockIndex++
		}
		if inst.Flags.IsJumpTarget() && i != 0 {
			fmt.Fprintf(&buf, "%s:\n", index2label[i])
			fmt.Fprintf(&buf, "  b%d %s\n", blockIndex, irfmt.Sprint(st, inst))
		} else {
			line := fmt.Sprintf("  b%d %s\n", blockIndex, irfmt.Sprint(st, inst))
			m := branchArgRE.FindStringSubmatch(line)
			if m != nil {
				index, err := strconv.Atoi(m[1])
				if err != nil {
					panic(err)
				}
				line = strings.Replace(line, m[0], index2label[index], 1)
			}
			buf.WriteString(line)
		}
	}

	return buf.String()
}
