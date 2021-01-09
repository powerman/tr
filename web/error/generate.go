package main

//go:generate rm -f main_wasm.go
//go:generate gobin -m -run github.com/vugu/vugu/cmd/vugugen -r -skip-go-mod
