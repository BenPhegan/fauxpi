package main

import (
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"net/http"
)

type StubResolver struct {
	FileChecker FileChecker
}

func (sr StubResolver) ReturnFileResponse() ResponseGenerator {
	return func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		var _, filename = constructFilename(r.Proto, r.Host, r.URL.Path, r.Method)
		var fileContent, _ = ioutil.ReadFile(filename)
		fileContentString := string(fileContent[:])
		statusCode := resolveStatusCode(fileContentString)

		return r, goproxy.NewResponse(r, "application/json", statusCode, fileContentString)
	}
}

func (sr StubResolver) CheckFilesystemForRequest() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		var _, filename = constructFilename(req.Proto, req.Host, req.URL.Path, req.Method)
		if _, err := sr.FileChecker(filename); err == nil {
			return true
		}
		return false
	}
}
