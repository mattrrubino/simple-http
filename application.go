package main

import (
	"fmt"
	"os"
	"strings"
)

func getFile(filepath string) (*os.File, int64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, 0, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	size := fileInfo.Size()

	return file, size, nil
}

func get404Response() (*HttpResponse, error) {
	var sb strings.Builder

	sb.WriteString("<html>\n")
	sb.WriteString("\t<head>\n")
	sb.WriteString("\t\t<title>Simple HTTP Server</title>\n")
	sb.WriteString("\t</head\n")
	sb.WriteString("\t<body>\n")
	sb.WriteString("\t\t404 Not Found\n")
	sb.WriteString("\t</body>\n")
	sb.WriteString("</html>\n")

	body := sb.String()
	bodyReader := strings.NewReader(body)

	contentLength := fmt.Sprint(len(body))
	headers := map[string]string{
		"Content-Type":   "text/html",
		"Content-Length": contentLength,
	}

	response := HttpResponse{"HTTP/1.1", 404, headers, bodyReader}
	return &response, nil
}

func listDirectory(filepath string) []string {
	files, err := os.ReadDir(filepath)
	if err != nil {
		return []string{}
	}

	filenames := make([]string, len(files)+1)
	filenames[0] = ".."
	for i, file := range files {
		filenames[i+1] = file.Name()
	}

	return filenames
}

func getDirectoryResponse(filepath string) (*HttpResponse, error) {
	var sb strings.Builder

	sb.WriteString("<html>\n")
	sb.WriteString("\t<head>\n")
	sb.WriteString("\t\t<title>Simple HTTP Server</title>\n")
	sb.WriteString("\t</head\n")
	sb.WriteString("\t<body>\n")

	for _, filename := range listDirectory(filepath) {
		fullPath := filepath + "/" + filename
		link := fmt.Sprintf("\t\t<a href=\"%v\">%v</a><br>\n", fullPath, filename)
		sb.WriteString(link)
	}

	sb.WriteString("\t</body>\n")
	sb.WriteString("</html>\n")

	body := sb.String()
	bodyReader := strings.NewReader(body)

	contentLength := fmt.Sprint(len(body))
	headers := map[string]string{
		"Content-Type":   "text/html",
		"Content-Length": contentLength,
	}

	response := HttpResponse{"HTTP/1.1", 200, headers, bodyReader}
	return &response, nil
}

func getFileResponse(filepath string) (*HttpResponse, error) {
	file, fileSize, err := getFile(filepath)
	if err != nil {
		return nil, err
	}

	contentLength := fmt.Sprint(fileSize)
	headers := map[string]string{
		"Connection":          "keep-alive",
		"Content-Disposition": "attachment",
		"Content-Type":        "application/octet-stream",
		"Content-Length":      contentLength,
	}

	response := HttpResponse{"HTTP/1.1", 200, headers, file}
	return &response, nil
}

func requestPathToFilePath(requestPath string) string {
	filepath := strings.Trim(requestPath, "/")

	// Empty path should map to present directory
	if filepath == "" {
		filepath = "."
	}

	return filepath
}

func getHttpResponse(request *HttpRequest) (*HttpResponse, error) {
	filepath := requestPathToFilePath(request.path)

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return get404Response()
	}

	if fileInfo.IsDir() {
		return getDirectoryResponse(filepath)
	} else {
		return getFileResponse(filepath)
	}
}
