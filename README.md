# go-jdk

![Logo](docs/logo_small.png)

[go-jdk](https://github.com/quasilyte/go-jdk) is [OpenJDK](https://ru.wikipedia.org/wiki/OpenJDK)-like implementation
written in Go with a goal to deliver a great embeddable [JVM](https://en.wikipedia.org/wiki/Java_virtual_machine) for
[Go](http://golang.org/) applications without [CGo](https://golang.org/cmd/cgo/).

Key features:

* JVM bytecode converted to [register-based form](https://www.usenix.org/legacy/events%2Fvee05%2Ffull_papers/p153-yunhe.pdf)
* Loaded code is JIT-compiled right away, no run-time tracing involved
* Cheap `Go->JVM` calls
* Cheap `JVM->Go` calls
* `native` Java methods can be written in Go
