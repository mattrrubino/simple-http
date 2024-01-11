# Simple HTTP Server

This repository contains the code for a simple HTTP server which serves static files from a directory. This server should NOT be used in production as it has not been tested. I built this as an exercise to learn Go.

## Usage

To run the server locally, execute the following commands:

1. `git clone git@github.com:mattrrubino/simple-http.git`
2. `cd simple-http`
3. `go build -o http`

Then, enter the directory you would like to serve in a terminal and run the executable. Note
that if you are on Windows, you should instead execute `go build -o http.exe` for step 3, and
you should run `start {path to http.exe}` from the directory you would like to serve.
The server runs on port 8080 by default, and this can be modified using the `-p` flag.
