package main

import (
	"errors"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
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

func Test_UrlChecker(t *testing.T) {
	var resolver = new(StubResolver)
	resolver.FileChecker = MockFileResolver
	var request = new(http.Request)
	request.RequestURI = "http://google.com/"
	var requestUrl, _ = url.Parse("http://google.com/")
	request.URL = requestUrl
	request.URL.Path = ""
	var context = new(goproxy.ProxyCtx)
	proxy := ReturnMockProxy()
	if resolver.CheckFilesystemForRequest().HandleReq(request, context) == true {
		t.Log("Passed")
	} else {
		t.Error("Failed")
	}

}

func Test_HostFileNameCreation(t *testing.T) {
	var hostfilename, filename = constructFilename("HTTP/1.1", "www.google.com", "", "GET")
	if hostfilename != "http/www.google.com/index.get.json" {
		t.Failed()
	}
}

func MockFileResolver(name string) (fi os.FileInfo, err error) {
	if name == "http/www.google.com/index.get.json" {
		return new(MockFileInfo), nil
	} else {
		return nil, errors.New("File not found")
	}
}

type MockFileInfo struct{}

func (m MockFileInfo) Name() string {
	return ""
}
func (m MockFileInfo) Size() int64 {
	return 1
}
func (m MockFileInfo) Mode() os.FileMode {
	return 1
}
func (m MockFileInfo) ModTime() time.Time {
	time := new(time.Time)
	return *time
}
func (m MockFileInfo) IsDir() bool {
	return false
}
func (m MockFileInfo) Sys() interface{} {
	return nil
}

func ReturnMockProxy() goproxy.ProxyHttpServer {
	proxy := goproxy.ProxyHttpServer{
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	return proxy
}
