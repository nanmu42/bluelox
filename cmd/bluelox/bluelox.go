package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/nanmu42/bluelox/interpreter"

	"github.com/nanmu42/bluelox/lox"
)

func main() {
	var (
		err      error
		exitCode int
	)
	defer func() {
		if err != nil {
			fmt.Println(err)
			os.Exit(exitCode)
		}
	}()

	if len(os.Args) > 2 {
		fmt.Println("Usage: bluelox [script]")
		exitCode = 64
		return
	}

	runner := lox.NewLox()
	if len(os.Args) == 1 {
		err = runner.RunPrompt()
		if err != nil {
			err = fmt.Errorf("running prompt: %w", err)
			return
		}
		exitCode = 65
		return
	}

	err = runner.RunFile(os.Args[1])
	if err != nil {
		err = fmt.Errorf("running script file: %w", err)

		var runtimeErr *interpreter.RuntimeError
		if errors.As(err, &runtimeErr) {
			exitCode = 70
		} else {
			exitCode = 65
		}

		return
	}
}
