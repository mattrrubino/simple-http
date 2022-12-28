/*
Http serves static files using HTTP.

Files are served from the directory from which the program is executed.
The server runs on port 8080 by default. This can be changed using flags.

Usage:

	http [flags]

The flags are:

	-p

		Set the port number of the HTTP server.
*/
package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("Starting HTTP server...")

	port := flag.Int("p", 8080, "Set the port number of the HTTP server.")
	flag.Parse()
	portString := fmt.Sprint(*port)

	if err := StartHttpServer("0.0.0.0", portString); err != nil {
		fmt.Printf("HTTP server error: %v\n", err)
		return
	}

	fmt.Println("Stopping HTTP server...")
}
