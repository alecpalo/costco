package main

import "costco/internal/costco"

func main() {
	r := costco.Init()
	r.RegisterEndpoints()
	r.Start()
}
