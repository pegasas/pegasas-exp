# Golang

## Compile

![gccgo_structure.png](./gccgo_structure.png)

Gccgo 作为 GCC（GNU 编译器集合）的一部分进行分发。 GCC 支持不同语言的几种前端。 gccgo 是连接到 GCC 后端的 Go 前端。 Go 前端与 GCC 项目是分开的，旨在能够连接到其他编译器后端，但目前仅支持 GCC。
与 gc 相比，gccgo 编译代码的速度较慢，但支持更强大的优化，因此 gccgo 构建的受 CPU 约束的程序通常运行速度更快。多年来，在 GCC 中实现的所有优化都是可用的，包括内联，循环优化，向量化，指令调度等。尽管它并不总是能产生更好的代码，但在某些情况下，使用 gccgo 编译的程序可以运行 30％更快。
gc 编译器仅支持最受欢迎的处理器：x86（32 位和 64 位）和 ARM。但是，Gccgo 支持 GCC 支持的所有处理器。并非所有这些处理器都经过了 gccgo 的全面测试，但是许多处理器都进行了测试，包括 x86（32 位和 64 位），SPARC，MIPS，PowerPC 甚至 Alpha。 Gccgo 还已经在 gc 编译器不支持的操作系统（尤其是 Solaris）上进行了测试。
Gccgo 提供了标准的，完整的 Go 库。 Go 运行时的许多核心功能在 gccgo 和 gc 中都相同，包括 goroutine 调度程序，通道，内存分配器和垃圾收集器。 Gccgo 和 gc 编译器一样，支持拆分 goroutine 堆栈，但目前仅在 x86（32 位或 64 位）上，并且仅在使用 gold 链接器时（在其他处理器上，每个 goroutine 将具有较大的堆栈，以及一系列深层次的函数调用可能会运行到堆栈末尾并使程序崩溃）。