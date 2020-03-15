# go-jdk

![Logo](docs/logo_small.png)

[go-jdk](https://github.com/quasilyte/go-jdk) is [OpenJDK](https://ru.wikipedia.org/wiki/OpenJDK)-like implementation
written in Go with a goal to deliver a great embeddable [JVM](https://en.wikipedia.org/wiki/Java_virtual_machine) for
[Go](http://golang.org/) applications without [CGo](https://golang.org/cmd/cgo/).

Key features:

* JVM bytecode converted to [register-based form](https://www.usenix.org/legacy/events%2Fvee05%2Ffull_papers/p153-yunhe.pdf)
* Loaded code is JIT-compiled right away, no run-time tracing involved
* Efficient `Go->JVM` calls
* Efficient `JVM->Go` calls
* `native` Java methods can be written in Go

> Note: this project is in its early state.

```bash
# Run Java class method (main or any other static method):
go-jdk run -class Foo.class -method helloWorld

# Disassemble Java class file with go-jdk:
go-jdk javap Foo.class

# Print IR representation instead of JVM bytecode:
go-jdk javap -format=ir Foo.class

# Print Java class dependencies:
go-jdk jdeps Foo.class
```
