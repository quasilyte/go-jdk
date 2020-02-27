# go-jdk

![Logo](docs/logo_small.png)

## Architecture overview

Main packages:

* [`jclass`](jclass) decodes [Java class files](https://en.wikipedia.org/wiki/Java_class_file)
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

Utility packages:

* [`bytecode`](bytecode) describes the [Java bytecode](https://en.wikipedia.org/wiki/Java_bytecode) opcodes
* [`irfmt`](irfmt) converts IR instructions to pretty strings
* [`javap`](javap) pretty-prints Java class files
