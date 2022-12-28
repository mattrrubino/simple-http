package main

import (
	"io"
	"strings"
)

// An HttpRequest represents a stream of bytes
// sent over a TCP connection conforming to HTTP.
type HttpRequest struct {
	version     string
	method      string
	path        string
	queryString map[string]string
	headers     map[string]string
}

// An HttpResponseBody implements a reader which is
// used to acquire the byte stream of the body.
type HttpResponseBody interface {
	io.Reader
}

// An HttpResponse represents a stream of bytes that will
// be sent over a TCP connection conforming to HTTP.
type HttpResponse struct {
	version string
	code    int
	headers map[string]string
	body    HttpResponseBody
}

// An HttpError describes an error that
// occurred while handling an HTTP request
type HttpError struct {
	message string
}

// Error prints the error message of an HttpError object.
func (error HttpError) Error() string {
	return error.message
}

// Map containing all valid methods for HTTP requests.
var validHttpMethods = map[string]bool{
	"GET":     true,
	"HEAD":    true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"CONNECT": true,
	"OPTIONS": true,
	"TRACE":   true,
	"PATCH":   true,
}

// parseHttpRequestHeaderStrings parses a list of strings from the
// header section of an HTTP request into a map from keys to values
// where each key-value pair is associated with one string from headerStrings
func parseHttpRequestHeaderStrings(headerStrings []string) map[string]string {
	headers := make(map[string]string)

	for _, headerString := range headerStrings {
		colonIndex := strings.Index(headerString, ":")

		if colonIndex == -1 {
			continue
		}

		key := strings.TrimSpace(headerString[:colonIndex])
		value := strings.TrimSpace(headerString[colonIndex+1:])

		if value == "" {
			continue
		}

		headers[key] = value
	}

	return headers
}

// parseHttpRequestTargetString parses the target in the start
// line of an HTTP request into its path and query strings.
func parseHttpRequestTargetString(target string) (string, map[string]string) {
	queryStrings := make(map[string]string)

	splitTargetString := strings.Split(target, "?")
	path := splitTargetString[0]
	if len(splitTargetString) < 2 {
		return path, queryStrings
	}

	query := splitTargetString[1]

	for _, queryString := range strings.Split(query, "&") {
		splitQueryString := strings.Split(queryString, "=")
		key := splitQueryString[0]
		value := splitQueryString[1]

		queryStrings[key] = value
	}

	return path, queryStrings
}

// parseHttpRequestBytes parses a slice of bytes associated with
// an HTTP request and converts it into an HttpRequest struct.
//
// It returns an error if the request is malformed.
func parseHttpRequestBytes(httpRequestBytes []byte) (*HttpRequest, error) {
	httpRequestString := string(httpRequestBytes)
	httpRequestLines := strings.Split(httpRequestString, "\n")

	if len(httpRequestLines) < 1 {
		return nil, HttpError{"Invalid HTTP request format."}
	}

	startLine := httpRequestLines[0]
	startLineElements := strings.Split(startLine, " ")

	if len(startLineElements) < 3 {
		return nil, HttpError{"Invalid HTTP start line."}
	}

	method := startLineElements[0]
	if !validHttpMethods[method] {
		return nil, HttpError{"Invalid HTTP method: " + method + "."}
	}

	path, queryStrings := parseHttpRequestTargetString(startLineElements[1])
	version := startLineElements[2]

	// TODO: Fix assumption that data is not sent
	headerStrings := httpRequestLines[1:]
	headers := parseHttpRequestHeaderStrings(headerStrings)

	return &HttpRequest{version, method, path, queryStrings, headers}, nil
}
