# go-jdk

![Logo](docs/logo_small.png)

`go-jdk` is an embeddable JVM for Go.

Key features:

* JVM bytecode converted to [register-based form](https://www.usenix.org/legacy/events%2Fvee05%2Ffull_papers/p153-yunhe.pdf)
* Eagerly JIT-complies all loaded code, no run-time tracing involved
* Cheap `Go->JVM` calls
* Cheap `JVM->Go` calls
* `native` Java methods can be written in Go

## Architecture overview

Package load path:

1. Package and all its dependencies are loaded from [Java class files](https://en.wikipedia.org/wiki/Java_class_file)
2. [Java bytecode](https://en.wikipedia.org/wiki/Java_bytecode) is then converted to IR
3. Generated IR gets optimized
4. JIT-compiler generates machine code from that IR
5. Loaded classes and packages are stored inside VM handle

### Main packages

* [`jclass`](jclass) decodes Java class files
* [`ir`](ir) describes our intermediate representation (IR)
* [`irgen`](irgen) converts bytecode into our IR
* [`iropt`](iropt) runs optimizations over IR
* [`jdeps`](jdeps) finds all package dependencies
* [`loader`](loader) fetches class files with all their dependencies
* [`jit/x64`](jit/x64) generates x86-64 machine code
* [`jit/compiler/x64`](jit/compiler/x64) turns IR into machine code
* [`vmdat`](vmdat) VM data structures that represent virtual machine state
* [`symbol`](symbol) defines index-like objects for efficient symbol referencing
* [`jruntime`](jruntime) implements loaded code runtime
* [`mmap`](mmap) wraps platform-dependent memory mapping code
* [`javatest`](javatest) runs tests written in Java against host VM

### Utility packages

* [`bytecode`](bytecode) describes the Java bytecode opcodes
* [`irfmt`](irfmt) converts IR instructions to pretty strings
* [`javap`](javap) pretty-prints Java class files
