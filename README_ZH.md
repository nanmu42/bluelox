# BlueLox

[English](https://godoc.org/github.com/nanmu42/bluelox) | **中文**

[![GoDoc](https://godoc.org/github.com/nanmu42/bluelox?status.svg)](https://godoc.org/github.com/nanmu42/bluelox)

BlueLox 是一个基于AST语法树的 Lox 解释器，使用 Golang 实现。

Lox 是 Robert Nystrom 设计的编程语言，在他精辟的书
[《Crafting Interpreters（手写两个解释器）》](https://craftinginterpreters.com/)中作为实现对象。
他在书里实现了一个Java的AST语法树解释器（jlox）和一个C的机器码解释器（clox），
附带循循善诱的解释，走心的手绘插图以及满满一袋关于早餐的比喻和笑话。

## Lox Playground

https://lox.nanmu.me/

一个基于WASM版本BlueLox的，在浏览器中运行的代码执行环境。

它可以作为你在学习和编程时的试错工具和参考实现。

## CLI

```bash
go install github.com/nanmu42/bluelox/cmd/bluelox@latest
```

命令行模式：

```bash
bluelox
```

执行某个文件：

```bash
bluelox script.lox
```

## 致谢

Lox 编程语言和 [Crafting Interpreters](https://craftinginterpreters.com/)
是 [Robert Nystrom](https://twitter.com/intent/user?screen_name=munificentbob) 的作品。

Lox Playground 从 [Go Playground](https://go.dev/play/) 汲取了很多点子、风格以及实现。

## License

Copyright © 2022 LI Zhennan

Released under Apache License 2.0.
