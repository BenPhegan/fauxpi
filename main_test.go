package main

import (
	"github.com/elazarl/goproxy"
	"net/http"
	"testing"
)

func Test_StubSource(t *testing.T) {
	var source = new(StubSource)
	if len(source.Stubs()) == 1 {
		t.Log("Passed")
	} else {
		t.Error("Failed")
	}
}

func Test_StubResolver(t *testing.T) {
	var source = new(StubSource)
	var resolver = new(StubResolver)
	resolver.StubSource = source
	var request = new(http.Request)
	request.RequestURI = "http://google.com/"
	var context = new(goproxy.ProxyCtx)
	if resolver.CheckUrlAgainstStub().HandleReq(request, context) == true {
		t.Log("Passed")
	} else {
		t.Error("Failed")
	}

}

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
var tests = []filenameTests{
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
}

func Test_MultipleHostFileNameCreation(t *testing.T) {
	for _, tt := range tests {
		var hostfilename, filename = constructFilename(tt.In.Protocol, tt.In.Host, tt.In.Path, tt.In.Method)
		if hostfilename != tt.Out.HostFilename {
			t.Error("Expected: " + tt.Out.HostFilename + " but got : " + hostfilename)
		}
		if filename != tt.Out.Filename {
			t.Error("Expected: " + tt.Out.Filename + " but got : " + filename)
		}
	}
}
