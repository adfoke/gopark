package main

import (
	"fmt"
	"gopark/config"
	"gopark/pkg/hello"
)

func main() {
	fmt.Println("Hello, world!")
	hello.SayHello()
	fmt.Println("Starting application:", config.AppConfig.AppName)
    fmt.Printf("Running on port: %d\n", config.AppConfig.Port)
    fmt.Printf("Debug mode: %v\n", config.AppConfig.Debug)
}