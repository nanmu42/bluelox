package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nanmu42/bluelox/version"

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

	version.SetSubName("cli")

	if len(os.Args) > 2 {
		fmt.Println("Usage: bluelox [script]")
		exitCode = 64
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	runner := lox.NewLox(os.Stdout)
	if len(os.Args) == 1 {
		err = runner.RunPrompt(ctx)
		if err != nil {
			err = fmt.Errorf("running prompt: %w", err)
			return
		}
		exitCode = 65
		return
	}

	err = runner.RunFile(ctx, os.Args[1])
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
