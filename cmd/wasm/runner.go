//go:build js && wasm

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"syscall/js"

	"github.com/nanmu42/bluelox/interpreter"
	"github.com/nanmu42/bluelox/lox"
)

var noopWriter = NoopWriter{}

type Runner struct {
	// protect status
	mu sync.Mutex

	// states
	running     bool
	cancelRun   context.CancelFunc
	interpreter *interpreter.Interpreter

	// TODO: implement us
	formatting bool
	cancelFmt  context.CancelFunc
}

func (r *Runner) Run(this js.Value, args []js.Value) (err error) {
	if argLength := len(args); argLength != 1 {
		return fmt.Errorf("want 1 arg, got %d", argLength)
	}
	script := args[0]
	if scriptType := script.Type(); scriptType != js.TypeString {
		return fmt.Errorf("want arg type %s, got %s", js.TypeString.String(), scriptType.String())
	}

	var lockReleased bool

	r.mu.Lock()
	defer func() {
		if !lockReleased {
			r.mu.Unlock()
		}
	}()
	if r.running {
		err = errors.New("already running, must stop at first")
		return
	}
	if r.formatting {
		err = errors.New("already formatting, must stop at first")
		return
	}

	stdout := newOutputWriter()
	runner := lox.NewLox(stdout)
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelRun = cancel
	r.running = true
	defer func() {
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
	}()

	// release the lock
	r.mu.Unlock()
	lockReleased = true

	err = runner.Run(ctx, []byte(script.String()))
	if errors.Is(err, context.Canceled) {
		// user cancel is ok
		err = nil
		return
	}
	if err != nil {
		return fmt.Errorf("loxrun: %s", err.Error())
	}

	return nil
}

func (r *Runner) Stop(this js.Value, args []js.Value) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		r.interpreter.ChangeStdoutTo(noopWriter)
		r.cancelRun()
	}
	if r.formatting {
		// TODO: deal with fmt
	}

	return nil
}

func (r *Runner) Fmt(this js.Value, args []js.Value) (err error) {
	// TODO: implement me
	return errors.New("loxfmt: not implemented yet")
}

type OutputWriter struct {
	written bool
	output  js.Value
}

func newOutputWriter() *OutputWriter {
	return &OutputWriter{
		written: false,
		output:  js.Global().Get("writeOutput"),
	}
}

// Write message to JS writer, schema:
// {
//     Kind: 'string', // 'start', 'stdout', 'stderr', 'end'
//     Body: 'string'  // content of write or end status message
// }
func (o *OutputWriter) Write(p []byte) (n int, err error) {
	if !o.written {
		o.output.Invoke(map[string]interface{}{
			"Kind": "start",
		})
		o.written = true
	}

	o.output.Invoke(map[string]interface{}{
		"Kind": "stdout",
		"Body": string(p),
	})

	return len(p), nil
}

// Close see write for schema
func (o *OutputWriter) Close() error {
	o.output.Invoke(map[string]interface{}{
		"Kind": "stdout",
		"Body": "",
	})

	return nil
}

type NoopWriter struct{}

func (w NoopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
