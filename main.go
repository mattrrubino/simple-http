package main

import "fmt"

func main() {
	fmt.Println("Starting HTTP server...")

	if err := StartHttpServer("0.0.0.0", "8080"); err != nil {
		fmt.Printf("HTTP server error: %v\n", err)
		return
	}

	fmt.Println("Stopping HTTP server...")
}
