package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	log.SetFlags(0)

	_, err := runCommand("go", "tool", "asm", "asmtest.s")
	if err != nil {
		log.Fatalf("asm: %v", err)
	}
	out, err := runCommand("go", "tool", "objdump", "asmtest.o")
	if err != nil {
		log.Fatalf("objdump: %v", err)
	}
	data, err := ioutil.ReadFile("asmtest.s")
	if err != nil {
		log.Fatalf("read file: %v", err)
	}
	sourceLines := strings.Split(string(data), "\n")

	type asmLine struct {
		LineNum int
		loc     string
		Enc     string
		Asm     string
		Comment string
	}
	type asmFunc struct {
		Name  string
		Lines []asmLine
	}

	var funcs []asmFunc
	for _, funcString := range strings.Split(out, "\n\n") {
		lines := strings.Split(funcString, "\n")
		var fn asmFunc
		fn.Name = strings.TrimSuffix(strings.Fields(lines[0])[1], "(SB)")
		for _, l := range lines[1:] {
			if l == "" {
				continue
			}
			fields := strings.Fields(l)
			loc := fields[0]
			enc := fields[2]
			var lineNum int
			fmt.Sscanf(loc, "asmtest.s:%d", &lineNum)
			src := sourceLines[lineNum-1]
			parts := strings.Split(src, "//")
			asm := parts[0]
			var comment string
			if len(parts) > 1 {
				comment = parts[1]
			}
			fn.Lines = append(fn.Lines, asmLine{
				LineNum: lineNum,
				loc:     loc,
				Enc:     enc,
				Asm:     strings.TrimSpace(asm),
				Comment: comment,
			})
		}
		lastLine := fn.Lines[len(fn.Lines)-1]
		if lastLine.Asm != "RET" {
			log.Fatalf("%s last line is not RET", fn.Name)
		}
		// Drop the last RET from the lines list.
		fn.Lines = fn.Lines[:len(fn.Lines)-1]
		funcs = append(funcs, fn)
	}

	tmpl := template.Must(template.New(`testcase`).Parse(`
{
  name: "{{$.Name}}",
  want: []expected{
    {{- range $.Lines}}
      { {{.LineNum}}, "{{.Asm}}", "{{.Enc}}" },
    {{- end}}
  },
  run: func(asm *Assembler) {
    {{- range $.Lines}} {{.Comment}}; {{end}}
  },
},
`))
	for _, fn := range funcs {
		err := tmpl.Execute(os.Stdout, fn)
		if err != nil {
			log.Fatalf("run %s template: %v", fn.Name, err)
		}
	}
}

/*
main.asmFunc{
  name:"jump3",
  lines:[]main.asmLine{
    main.asmLine{
      loc:"asmtest.s:15",
      enc:"90",
      asm:"NOP1",
      comment:" asm.Label(0); asm.Nop(1)"
    },
    main.asmLine{
      loc:"asmtest.s:16",
      enc:"ebfd",
      asm:"JMP backward",
      comment:" asm.Jmp(0)"
    },
    main.asmLine{
      loc:"asmtest.s:17",
      enc:"c3",
      asm:"RET",
      comment:"",
    }
  }
}

*/

func runCommand(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return string(out), nil
}
