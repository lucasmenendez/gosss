//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/lucasmenendez/gosss"
)

const (
	// js names
	jsClassName     = "GoSSS"
	jsHideMethod    = "hide"
	jsRecoverMethod = "recover"
	// required number of arguments for each method
	hidedNArgs    = 3 // msg, nshares, minshares
	recoverdNArgs = 1 // shares
)

func wasmError(err error) js.Value {
	return js.Global().Call("throwError", js.ValueOf(err.Error()))
}

func main() {
	gosssClass := js.ValueOf(map[string]interface{}{})
	gosssClass.Set(jsHideMethod, js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		if len(p) != hidedNArgs {
			return wasmError(fmt.Errorf("invalid number of arguments"))
		}
		msg := p[0].String()
		nshares := p[1].Int()
		minshares := p[2].Int()
		// hide the message
		shares, err := gosss.HideMessage([]byte(msg), &gosss.Config{
			Shares: nshares,
			Min:    minshares,
		})
		if err != nil {
			return wasmError(err)
		}
		jsonShares, err := json.Marshal(shares)
		if err != nil {
			return wasmError(err)
		}
		return string(jsonShares)
	}))

	gosssClass.Set(jsRecoverMethod, js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		if len(p) != recoverdNArgs {
			return wasmError(fmt.Errorf("invalid number of arguments"))
		}
		// recover the shares from the json string
		var shares []string
		if err := json.Unmarshal([]byte(p[0].String()), &shares); err != nil {
			return wasmError(err)
		}
		// recover the message
		msg, err := gosss.RecoverMessage(shares, nil)
		if err != nil {
			return wasmError(err)
		}
		return string(msg)
	}))

	js.Global().Set(jsClassName, gosssClass)
	// keep the program running forever
	select {}
}
