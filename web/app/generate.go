package main

// Remove main_wasm.go to update it in case of vugu upgrade.
//go:generate rm -f main_wasm.go
//go:generate gobin -m -run github.com/vugu/vugu/cmd/vugugen -skip-go-mod
