package javatest

import (
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/quasilyte/go-jdk/irgen"
	"github.com/quasilyte/go-jdk/jit"
	"github.com/quasilyte/go-jdk/jruntime"
	"github.com/quasilyte/go-jdk/loader"
)

var testsDebug = os.Getenv("DEBUG") == "true"

var tests = []*testParams{
	{Pkg: "intvalues", Input: 400},
	{Pkg: "longvalues", Input: 400},
	{Pkg: "scopes"},
	{Pkg: "arith1", Input: 100},
	{Pkg: "staticcall1"},
	{Pkg: "loops1"},
}

func TestMain(m *testing.M) {
	if !hasCommand("java") {
		log.Println("skip: missing java")
		return
	}
	if !hasCommand("javac") {
		log.Println("skip: missing javac")
		return
	}

	fillTestDefaults(tests)
	generateJavaMain(tests)
	compileJava()
	code := m.Run()
	os.Remove(filepath.Join("testdata", "Main.java"))
	os.Exit(code)
}

func TestJava(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Pkg, func(t *testing.T) {
			runTest(t, test)
		})
	}
}

type testParams struct {
	Pkg string

	EntryClass  string
	EntryMethod string

	Input int32
}

func fillTestDefaults(tests []*testParams) {
	for _, test := range tests {
		if test.EntryClass == "" {
			test.EntryClass = "Test"
		}
		if test.EntryMethod == "" {
			test.EntryMethod = "run"
		}
	}
}

func runTest(t *testing.T, params *testParams) {
	vm, err := jruntime.OpenVM(runtime.GOARCH)
	if err != nil {
		t.Fatalf("open VM: %v", err)
	}
	defer vm.Close()

	vm.State.BindGoFunc("testutil/T.printInt", golibPrintInt)
	vm.State.BindGoFunc("testutil/T.printLong", golibPrintLong)

	absTestdata, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatalf("abs(testdata): %v", err)
	}
	packages, err := loader.LoadPackage(&vm.State, params.Pkg, &loader.Config{
		ClassPath: []string{
			absTestdata,
			filepath.Join(absTestdata, "_golib"),
		},
	})
	if err != nil {
		t.Fatalf("load package: %v", err)
	}
	if err := irgen.Generate(&vm.State, packages); err != nil {
		t.Fatalf("irgen: %v", err)
	}
	ctx := jit.Context{
		Mmap:  &vm.Mmap,
		State: &vm.State,
	}
	if err := vm.Compiler.Compile(ctx, packages); err != nil {
		t.Fatalf("compile: %v", err)
	}

	class := packages[0].Out.FindClass(params.EntryClass)
	if class == nil {
		t.Fatalf("entry class %s not found", params.EntryClass)
	}
	method := class.FindMethod(params.EntryMethod, "")
	if method == nil {
		t.Fatalf("entry method %s not found", params.EntryMethod)
	}
	if method.Descriptor != "(I)V" {
		t.Fatalf("entry method signature should be: void %s()", params.EntryMethod)
	}

	golibOutput.Reset()
	env := jruntime.NewEnv(vm, &jruntime.EnvConfig{})
	env.IntArg(0, int64(params.Input))
	if _, err := env.IntCall(method); err != nil {
		t.Fatalf("call error: %v", err)
	}
	have := golibOutput.String()
	want := runJava(t, params)
	if testsDebug {
		t.Logf("Go output:\n%s", have)
	}
	if have != want {
		t.Errorf("output mismatch:\nhave:\n%s\nwant:\n%s", have, want)
	}
}

func runJava(t *testing.T, params *testParams) string {
	args := []string{
		"-cp", "testdata:testdata/_javalib",
		"Main",
		params.Pkg,
	}
	out, err := exec.Command("java", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("java: %v: %s", err, out)
	}
	return string(out)
}

func generateJavaMain(tests []*testParams) {
	tmpl := template.Must(template.New(`Main`).Parse(`
// Generated automatically by java_test.go.
// This entry point is used by a host Java implementation.
class Main {
    public static void main(String args[]) {
        switch (args[0]) {
        {{- range .}}
        case "{{.Pkg}}":
            {{.Pkg}}.{{.EntryClass}}.{{.EntryMethod}}({{.Input}});
            return;
        {{- end}}
        default:
            System.out.println("unknown package: " + args[0]);
        }
    }
}
`))
	f, err := os.Create(filepath.Join("testdata", "Main.java"))
	if err != nil {
		log.Fatalf("create file: %v", err)
	}
	defer f.Close()
	if err := tmpl.Execute(f, tests); err != nil {
		log.Fatalf("execute template: %v", err)
	}
}

func compileJava() {
	// Compilation of Main.java will create class files for all tests.
	{
		args := []string{
			"-cp", "testdata:testdata/_javalib",
			"testdata/Main.java",
		}
		out, err := exec.Command("javac", args...).CombinedOutput()
		if err != nil {
			log.Fatalf("javac: %v: %s", err, out)
		}
	}

	// Now we only need to compile _golib classes.
	{
		args := []string{
			"-cp", "testdata/_golib",
			"testdata/_golib/testutil/T.java",
		}
		out, err := exec.Command("javac", args...).CombinedOutput()
		if err != nil {
			log.Fatalf("javac: %v: %s", err, out)
		}
	}
}

func hasCommand(name string) bool {
	err := exec.Command("/bin/sh", "-c", "command -v "+name).Run()
	return err == nil
}
