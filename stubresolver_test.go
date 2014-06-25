package main

import (
	"strconv"
	"testing"
)

type filenameParams struct {
	Protocol string
	Host     string
	Path     string
	Method   string
}

type filenameResults struct {
	HostFilename string
	Filename     string
}

type filenameTests struct {
	In  filenameParams
	Out filenameResults
}

//HTTP/1.1,google.com,/test,GET
var fileTests = []filenameTests{
	{
		filenameParams{"HTTP/1.1", "www.google.com", "/", "GET"},
		filenameResults{"http/www.google.com/index.get.json", "index.get.json"},
	},
	{
		filenameParams{"HTTP/1.1", "www.google.com", "/search", "GET"},
		filenameResults{"http/www.google.com/_search.get.json", "_search.get.json"},
	},
	{
		filenameParams{"HTTP/1.1", "www.google.com", "/comments/", "GET"},
		filenameResults{"http/www.google.com/comments/index.get.json", "comments/index.get.json"},
	},
	{
		filenameParams{"HTTP/1.1", "www.google.com", "/comments/7", "GET"},
		filenameResults{"http/www.google.com/comments/any.get.json", "comments/any.get.json"},
	},
	{
		filenameParams{"HTTP/1.1", "www.google.com", "/7", "GET"},
		filenameResults{"http/www.google.com/any.get.json", "any.get.json"},
	},
}

func Test_MultipleHostFileNameCreation(t *testing.T) {
	for _, tt := range fileTests {
		var hostfilename, filename = constructFilename(tt.In.Protocol, tt.In.Host, tt.In.Path, tt.In.Method)
		if hostfilename != tt.Out.HostFilename {
			t.Error("Expected: " + tt.Out.HostFilename + " but got : " + hostfilename)
		}
		if filename != tt.Out.Filename {
			t.Error("Expected: " + tt.Out.Filename + " but got : " + filename)
		}
	}
}

var statusCodeTests = []struct {
	text       string
	statusCode int
}{
	{
		text:       "//! statusCode: 201 <html> <body>Created something successfully! Happy!</body></html>",
		statusCode: 201,
	},
	{
		text:       "//! statusCode: 500 <html> <body>BOOM</body></html>",
		statusCode: 500,
	},
	{
		text:       "<html> <body>BOOM</body></html>",
		statusCode: 200,
	},
}

func Test_CanCreateCustomStatusCodes(t *testing.T) {
	for _, tt := range statusCodeTests {
		response := resolveStatusCode(tt.text)
		if response != tt.statusCode {
			t.Error("Expected: " + strconv.Itoa(tt.statusCode) + " but got : " + strconv.Itoa(response))
		}
	}
}
