package main

import "costco/internal/registry"

func main() {
	r := registry.Init()
	r.RegisterEndpoints()
	r.Start()
}
