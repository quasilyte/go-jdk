// Package javap is a very simplistic Java class files dumping utility.
//
// Mostly intended for testing and debugging.
package javap

import (
	"fmt"
	"io"
	"strings"

	"github.com/quasilyte/go-jdk/bytecode"
	"github.com/quasilyte/go-jdk/jclass"
)

// Fprint pretty-prints c to a given writer.
func Fprint(w io.Writer, c *jclass.File) {
	p := printer{w: w, c: c}

	className := c.ThisClassName
	p.write("class file for %q\n", className)
	p.write("version: %d.%d\n", c.Ver.Major, c.Ver.Minor)

	p.write("constants:\n")
	for i, c := range c.Consts[1:] {
		s := strings.ReplaceAll(fmt.Sprintf("%#v", c), "&jclass.", "&")
		p.write("  %-3d %s\n", i+1, s)
	}

	p.write("methods:\n")
	for _, m := range c.Methods {
		p.printMethod(m)
	}
}

type printer struct {
	w io.Writer
	c *jclass.File
}

func (p *printer) write(format string, args ...interface{}) {
	fmt.Fprintf(p.w, format, args...)
}

func (p *printer) printMethod(m jclass.Method) {
	sig := jclass.MethodDescriptor(m.Descriptor).SignatureString(m.Name)
	if m.AccessFlags.IsNative() {
		p.write("  native method %s\n", sig)
		return
	}

	p.write("  method %s:\n", sig)
	codeAttr := findAttr(p.c, m.Attrs, "Code").(jclass.CodeAttribute)
	p.write("    max_locals=%d max_stack=%d\n",
		codeAttr.MaxLocals, codeAttr.MaxStack)
	p.write("    bytecode (size=%d):\n", len(codeAttr.Code))

	pc := 0
	code := codeAttr.Code
	for pc < len(code) {
		op := bytecode.Op(code[pc])
		width := int(bytecode.OpWidth[op])
		opbytes := code[pc : pc+width]
		fmt.Printf("      %3x %-18s %x\n", pc, op.String(), opbytes)
		pc += width
	}
	frameTab, ok := findAttr(p.c, codeAttr.Attrs, "StackMapTable").(jclass.StackMapTableAttribute)
	if ok && len(frameTab.Frames) != 0 {
		fmt.Println("      ---- (stack map) ----")
		for _, frame := range frameTab.Frames {
			fmt.Printf("      %3x depth=%d\n", frame.Offset, frame.StackDepth)
		}
	}

	p.write("\n")
}

func findAttr(c *jclass.File, attrs []jclass.Attribute, name string) jclass.Attribute {
	for _, attr := range attrs {
		switch attr := attr.(type) {
		case jclass.StackMapTableAttribute:
			if name == "StackMapTable" {
				return attr
			}
		case jclass.CodeAttribute:
			if name == "Code" {
				return attr
			}
		case jclass.RawAttribute:
			attrName := c.Consts[attr.NameIndex].(*jclass.Utf8Const).Value
			if name == attrName {
				return attr
			}
		}
	}
	return nil
}
