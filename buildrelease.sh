#!/bin/bash

rm ./ui/lib.wasm
GOOS=windows GOARCH=amd64 go build -o wowsim.exe web.go
GOOS=darwin GOARCH=amd64 go build -o wowsim-amd64-darwin web.go
GOOS=linux go build -o wowsim-amd64-linux web.go
cd ui && GOARCH=wasm GOOS=js go build -o lib.wasm main_wasm.go