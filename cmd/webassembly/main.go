//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"syscall/js"

	"github.com/lucasmenendez/gosss"
)

const (
	// js names
	jsClassName     = "GoSSS"
	jsHideMethod    = "hide"
	jsRecoverMethod = "recover"
	jsMaxLenMethod  = "maxLength"
	// required number of arguments for each method
	hidedNArgs    = 3 // msg, nshares, minshares, prime(*)
	recoverdNArgs = 1 // shares, prime(*)
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
	gosssClass := js.ValueOf(map[string]interface{}{
		"defaultPrime": gosss.DefaultPrime.String(),
		"minShares":    gosss.MinShares,
		"minMinShares": gosss.MinMinShares,
	})
	gosssClass.Set(jsHideMethod, js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		if len(p) < hidedNArgs {
			return wasmResult(nil, fmt.Errorf("invalid number of arguments"))
		}
		msg := p[0].String()
		conf := &gosss.Config{
			Shares: p[1].Int(),
			Min:    p[2].Int(),
		}
		if len(p) > hidedNArgs {
			strPrime := p[3].String()
			var ok bool
			if conf.Prime, ok = new(big.Int).SetString(strPrime, 10); !ok {
				return wasmResult(nil, fmt.Errorf("invalid prime number"))
			}
		}
		// hide the message
		shares, err := gosss.HideMessage([]byte(msg), conf)
		if err != nil {
			return wasmResult(nil, err)
		}
		return wasmResult(shares, nil)
	}))

	gosssClass.Set(jsRecoverMethod, js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		if len(p) < recoverdNArgs {
			return wasmResult(nil, fmt.Errorf("invalid number of arguments"))
		}
		// recover the shares from the json string
		var shares []string
		if err := json.Unmarshal([]byte(p[0].String()), &shares); err != nil {
			return wasmResult(nil, err)
		}
		conf := &gosss.Config{}
		if len(p) > recoverdNArgs {
			strPrime := p[1].String()
			var ok bool
			if conf.Prime, ok = new(big.Int).SetString(strPrime, 10); !ok {
				return wasmResult(nil, fmt.Errorf("invalid prime number"))
			}
		}
		// recover the message
		msg, err := gosss.RecoverMessage(shares, conf)
		if err != nil {
			return wasmResult(nil, err)
		}
		return wasmResult(msg, nil)
	}))

	gosssClass.Set(jsMaxLenMethod, js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		conf := &gosss.Config{}
		if len(p) > 0 {
			strPrime := p[0].String()
			var ok bool
			if conf.Prime, ok = new(big.Int).SetString(strPrime, 10); !ok {
				return wasmResult(nil, fmt.Errorf("invalid prime number"))
			}
		}
		if err := conf.ValidPrime(); err != nil {
			return wasmResult(nil, err)
		}
		return wasmResult(conf.MaxMessageLen(), nil)
	}))

	js.Global().Set(jsClassName, gosssClass)
	// keep the program running forever
	select {}
}
