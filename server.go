package main

import (
	"fmt"
	"io"
	"net"
)

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
