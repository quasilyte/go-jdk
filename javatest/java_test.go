package javatest

import (
	"fmt"
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
	"github.com/quasilyte/go-jdk/vmdat"
)

var testsDebug = os.Getenv("DEBUG") == "true"

var tests = []*testParams{
	{Pkg: "intvalues", Input: 400},
	{Pkg: "longvalues", Input: -400},
	{Pkg: "scopes"},
	{Pkg: "arith1", Input: 100},
	{Pkg: "gocall1", Input: -100},
	{Pkg: "gocall2"},
	{Pkg: "staticcall1"},
	{Pkg: "staticcall2"},
	{Pkg: "staticcall3"},
	{Pkg: "staticfinal"},
	{Pkg: "loops1"},
	{Pkg: "arrays1"},
	{Pkg: "arrays2"},
	{Pkg: "bubblesort"},
	{Pkg: "arrayreverse"},
	{Pkg: "eratosthenes", Input: 30},
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

func BenchmarkJava(b *testing.B) {
	benchmarks := []*benchParams{
		{
			Method: "nop",
		},

		{
			Method: "callGoNop",
			Ops:    5,
		},
	}
	runBenchmarks(b, benchmarks)
}

type benchParams struct {
	Name   string
	Method string
	Ops    int
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

func runBenchmarks(b *testing.B, benchmarks []*benchParams) {
	vm, err := jruntime.OpenVM(runtime.GOARCH)
	if err != nil {
		b.Fatalf("open VM: %v", err)
	}
	defer vm.Close()

	vm.State.BindGoFunc("benchutil/B.nop", golibNop)

	pkg, err := loadAndCompilePackage(vm, "bench")
	if err != nil {
		b.Fatal(err)
	}

	for _, params := range benchmarks {
		method, err := findMethod(pkg, "Bench", params.Method)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(params.Method, func(b *testing.B) {
			n := b.N
			if params.Ops > 1 && b.N != 1 {
				n /= params.Ops
			}
			env := jruntime.NewEnv(vm, &jruntime.EnvConfig{})
			for i := 0; i < n; i++ {
				env.IntCall(method)
			}
		})
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
	vm.State.BindGoFunc("testutil/T.printIntArray", golibPrintIntArray)
	vm.State.BindGoFunc("testutil/T.isub", golibIsub)
	vm.State.BindGoFunc("testutil/T.isub3", golibIsub3)
	vm.State.BindGoFunc("testutil/T.ii_l", golibII_L)
	vm.State.BindGoFunc("testutil/T.il_i", golibIL_I)
	vm.State.BindGoFunc("testutil/T.li_i", golibLI_I)
	vm.State.BindGoFunc("testutil/T.ilil_i", golibILIL_I)
	vm.State.BindGoFunc("testutil/T.GC", golibGC)

	pkg, err := loadAndCompilePackage(vm, params.Pkg)
	if err != nil {
		t.Fatal(err)
	}

	method, err := findMethod(pkg, params.EntryClass, params.EntryMethod)
	if err != nil {
		t.Fatal(err)
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
	// TODO: compile with 1 javac call instead of 2?

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
			"testdata/_golib/benchutil/B.java",
			"testdata/bench/Bench.java",
		}
		out, err := exec.Command("javac", args...).CombinedOutput()
		if err != nil {
			log.Fatalf("javac: %v: %s", err, out)
		}
	}
}

func findMethod(pkg *vmdat.Package, className, methodName string) (*vmdat.Method, error) {
	class := pkg.FindClass(className)
	if class == nil {
		return nil, fmt.Errorf("class %s not found in %s", className, pkg.Name)
	}
	method := class.FindMethod(methodName, "")
	if method == nil {
		return nil, fmt.Errorf("method %s not found in %s", methodName, class.Name)
	}
	return method, nil
}

func loadAndCompilePackage(vm *jruntime.VM, pkg string) (*vmdat.Package, error) {
	absTestdata, err := filepath.Abs("testdata")
	if err != nil {
		return nil, err
	}
	packages, err := loader.LoadPackage(&vm.State, pkg, &loader.Config{
		ClassPath: []string{
			absTestdata,
			filepath.Join(absTestdata, "_golib"),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("load %s: %v", pkg, err)
	}

	if err := irgen.Generate(&vm.State, packages); err != nil {
		return nil, fmt.Errorf("irgen: %v", err)
	}

	ctx := jit.Context{
		Mmap:  &vm.Mmap,
		State: &vm.State,
	}
	jruntime.BindFuncs(&ctx)
	if err := vm.Compiler.Compile(ctx, packages); err != nil {
		return nil, fmt.Errorf("compile: %v", err)
	}

	return packages[0].Out, nil
}

func hasCommand(name string) bool {
	err := exec.Command("/bin/sh", "-c", "command -v "+name).Run()
	return err == nil
}
