# BlueLox

**English** | [中文](https://github.com/nanmu42/bluelox/blob/master/README_ZH.md)

[![GoDoc](https://godoc.org/github.com/nanmu42/bluelox?status.svg)](https://godoc.org/github.com/nanmu42/bluelox)

BlueLox is a Tree-walking interpreter implemented in Golang for Lox.

Lox is a programing language by Robert Nystrom, 
introduced in his wonderful book [Crafting Interpreters](https://craftinginterpreters.com/), 
where he constructs a Java version interpreter(jlox) line by line, 
with detailed tutorial, brilliant illustrations and a full pack of jokes about breakfast.

## Lox Playground

https://lox.nanmu.me/

A web browser based Lox playground powered by WASM version of BlueLox.

You may find the Lox Playground helpful during your learning and coding as it may be your stage
for trial-and-error and implementation reference.

## CLI

```bash
go install github.com/nanmu42/bluelox/cmd/bluelox@latest
```

To use as a prompt:

```bash
bluelox
```

To run a script file:

```bash
bluelox script.lox
```

## Acknowledgement

Lox programing language and [Crafting Interpreters](https://craftinginterpreters.com/)
are works by [Robert Nystrom](https://twitter.com/intent/user?screen_name=munificentbob).

The Lox Playground borrows lots of ideas, code and styles from [Go Playground](https://go.dev/play/).

## License

Copyright © 2022 LI Zhennan

Released under Apache License 2.0.
