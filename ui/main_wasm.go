package main

import (
	"bytes"
	"encoding/json"
	"syscall/js"
	"github.com/lologarithm/wowsim/tbc"
)

func ApiCallJson(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		panic("API call requires a request!")
	}

	// Assumes input is a JSON object as a string
	requestStr := args[0].String()

	request := tbc.ApiRequest{}
	dec := json.NewDecoder(bytes.NewReader([]byte(requestStr)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&request); err != nil {
		panic("Error parsing request: " + err.Error() + "\nRequest: " + requestStr)
	}

	result := tbc.ApiCall(request)

	resultData, err := json.Marshal(result)
	if err != nil {
		panic("Error marshaling result: " + err.Error())
	}
	return string(resultData)
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("apiCall", js.FuncOf(ApiCallJson))
	js.Global().Call("wasmready")
	<-c
}
