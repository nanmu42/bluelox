//go:build js && wasm

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"syscall/js"

	"github.com/nanmu42/bluelox/lox"
)

var noopWriter = NoopWriter{}

type Runner struct {
	// protect status
	mu sync.Mutex

	// states
	running   bool
	cancelRun context.CancelFunc
	lox       *lox.Lox

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
		err = errors.New("stopping running program, wait for 'Program exited' and click Run to try again. Press F5 to refresh if it takes too long")
		return
	}
	if r.formatting {
		err = errors.New("stopping formatting program, please try again")
		return
	}

	stdout := newOutputWriter()
	r.lox = lox.NewLox(stdout)
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelRun = cancel
	r.running = true
	defer func() {
		_ = stdout.Close()

		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
	}()

	// release the lock
	r.mu.Unlock()
	lockReleased = true

	err = r.lox.Run(ctx, []byte(script.String()))
	if errors.Is(err, context.Canceled) {
		// user cancel is ok
		err = nil
		return
	}
	if err != nil {
		return fmt.Errorf("loxrun: %w", err)
	}

	return nil
}

func (r *Runner) Stop(this js.Value, args []js.Value) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		r.lox.ChangeStdoutTo(noopWriter)
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

func (o *OutputWriter) initWrite() {
	if !o.written {
		o.output.Invoke(map[string]interface{}{
			"Kind": "start",
		})
		o.written = true
	}
}

// Write message to JS writer, schema:
// {
//     Kind: 'string', // 'start', 'stdout', 'stderr', 'end'
//     Body: 'string'  // content of write or end status message
// }
func (o *OutputWriter) Write(p []byte) (n int, err error) {
	o.initWrite()

	o.output.Invoke(map[string]interface{}{
		"Kind": "stdout",
		"Body": string(p),
	})

	return len(p), nil
}

// Close see write for schema
func (o *OutputWriter) Close() error {
	o.initWrite()

	o.output.Invoke(map[string]interface{}{
		"Kind": "end",
		"Body": "",
	})

	return nil
}

type NoopWriter struct{}

func (w NoopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
