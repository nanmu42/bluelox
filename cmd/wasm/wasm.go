//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/nanmu42/bluelox/version"
)

var Promise = js.Global().Get("Promise")

func main() {
	version.SetSubName("wasm")

	runner := &Runner{}

	js.Global().Set("loxrun", asyncFuncOf(runner.Run))
	js.Global().Set("loxfmt", asyncFuncOf(runner.Fmt))
	js.Global().Set("loxstop", asyncFuncOf(runner.Stop))
	js.Global().Set("loxversion", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return version.FullNameWithBuildDate
	}))

	select {}
}

// asyncFuncOf avoids Go's js deadlock
//
// source: https://github.com/golang/go/issues/41310#issuecomment-725809881
func asyncFuncOf(fn func(this js.Value, args []js.Value) error) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(_ js.Value, promise []js.Value) interface{} {
			resolve := promise[0]
			reject := promise[1]

			go func() {
				err := fn(this, args)
				if err != nil {
					reject.Invoke(err.Error())
					return
				}

				resolve.Invoke()
			}()

			return nil
		})

		return Promise.New(handler)
	})
}
