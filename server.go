package main

import (
	"fmt"
	"io"
	"net"
)

// StartHttpServer starts an HTTP server on the specified ip and port.
//
// It returns an error if the server fails to bind to the specified address.
func StartHttpServer(ip, port string) error {
	address := ip + ":" + port
	ln, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	fmt.Printf("HTTP server running on %v\n", address)
	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Printf("Connection creation error: %v\n", err)
			continue
		}

		go handleTcpConnection(conn)
	}
}

// handleTcpConnection reads from and writes to a TCP connection following HTTP.
// It should be run as a goroutine to maximize server throughput.
func handleTcpConnection(conn net.Conn) {
	fmt.Printf("Opened TCP connection to %v\n", conn.RemoteAddr())
	defer conn.Close()

	httpRequest, err := getHttpRequest(conn)
	if err != nil {
		fmt.Printf("Error getting HTTP request: %v\n", err)
		return
	}

	httpResponse, err := getHttpResponse(httpRequest)
	if err != nil {
		fmt.Printf("Error getting HTTP response: %v\n", err)
		return
	}

	err = sendHttpResponse(conn, httpResponse)
	if err != nil {
		fmt.Printf("Error sending HTTP response: %v\n", err)
		return
	}
}

// getHttpRequest reads request bytes and transforms them
// into an HttpRequest struct.
//
// It returns an error if any read operations fail or if
// the request is malformed.
func getHttpRequest(conn net.Conn) (*HttpRequest, error) {
	httpRequestBytes := make([]byte, 1024)

	_, err := conn.Read(httpRequestBytes)
	if err != nil {
		return nil, err
	}

	httpRequest, err := parseHttpRequestBytes(httpRequestBytes)
	if err != nil {
		return nil, err
	}

	return httpRequest, nil
}

// sendHttpResponse transforms an HttpResponse struct into bytes
// and writes the data to the specified connection.
//
// It returns an error if any write operations fail.
func sendHttpResponse(conn net.Conn, httpResponse *HttpResponse) error {
	statusLine := fmt.Sprintf("%v %v\n", httpResponse.version, httpResponse.code)

	// Write status line
	_, err := conn.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	// Write headers
	for key, value := range httpResponse.headers {
		header := fmt.Sprintf("%v: %v\n", key, value)

		_, err = conn.Write([]byte(header))
		if err != nil {
			return err
		}
	}

	// Allow response without body
	if httpResponse.body == nil {
		return nil
	}

	// Write break for body
	_, err = conn.Write([]byte("\n"))
	if err != nil {
		return err
	}

	// Write body
	bodyData := make([]byte, 1024)
	for {
		n, err := httpResponse.body.Read(bodyData)

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		_, err = conn.Write(bodyData[:n])
	}
}
