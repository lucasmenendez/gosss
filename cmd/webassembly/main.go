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
	jsLimitsMethod  = "limits"
	// required number of arguments for each method
	hidedNArgs    = 3 // msg, nshares, minshares
	recoverdNArgs = 1 // shares
	limitsNArgs   = 1 // msg
)

func wasmResult(data interface{}, err error) js.Value {
	response := map[string]interface{}{}
	if data != nil {
		response["data"] = data
	}
	if err != nil {
		response["error"] = err.Error()
	}
	result, err := json.Marshal(response)
	if err != nil {
		return js.ValueOf(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}
	return js.ValueOf(string(result))
}

func main() {
	gosssClass := js.ValueOf(map[string]interface{}{})
	gosssClass.Set(jsHideMethod, js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		if len(p) != hidedNArgs {
			return wasmResult(nil, fmt.Errorf("invalid number of arguments"))
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
			return wasmResult(nil, err)
		}
		return wasmResult(shares, nil)
	}))

	gosssClass.Set(jsRecoverMethod, js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		if len(p) != recoverdNArgs {
			return wasmResult(nil, fmt.Errorf("invalid number of arguments"))
		}
		// recover the shares from the json string
		var shares []string
		if err := json.Unmarshal([]byte(p[0].String()), &shares); err != nil {
			return wasmResult(nil, err)
		}
		// recover the message
		msg, err := gosss.RecoverMessage(shares, nil)
		if err != nil {
			return wasmResult(nil, err)
		}
		return wasmResult(msg, nil)
	}))

	gosssClass.Set(jsLimitsMethod, js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		if len(p) != limitsNArgs {
			return wasmResult(nil, fmt.Errorf("invalid number of arguments"))
		}
		// recover the message
		minShares, maxShares, minMin, maxMin := gosss.ConfigLimits([]byte(p[0].String()))
		return wasmResult([]int{minShares, maxShares, minMin, maxMin}, nil)
	}))

	js.Global().Set(jsClassName, gosssClass)
	// keep the program running forever
	select {}
}
