package main

import (
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type StubResolver struct {
	FileChecker        FileChecker
	UseHostAndProtocol bool
	StubRoot           string
}

func (sr StubResolver) ReturnFileResponse() ResponseGenerator {
	return func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		var resolvedFilename = constructFilename(r.Proto, r.Host, r.URL.Path, r.Method, sr.UseHostAndProtocol, sr.StubRoot)

		var fileContent, _ = ioutil.ReadFile(resolvedFilename)
		fileContentString := string(fileContent[:])
		statusCode := resolveStatusCode(fileContentString)

		fileContentString = stripMetaData(fileContentString)

		return r, goproxy.NewResponse(r, "application/json", statusCode, fileContentString)
	}
}

func (sr StubResolver) CheckFilesystemForRequest() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		var filename = constructFilename(req.Proto, req.Host, req.URL.Path, req.Method, sr.UseHostAndProtocol, sr.StubRoot)
		if _, err := sr.FileChecker(filename); err == nil {
			return true
		}
		return false
	}
}

type ResponseGenerator func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)

func constructFilename(proto string, host string, reqPath string, method string, useHostAndProtocol bool, basePath string) string {
	verb := strings.ToLower(method)
	protocol := strings.ToLower(strings.Split(proto, "/")[0])
	urlpath := reqPath
	if strings.HasSuffix(urlpath, "/") {
		urlpath = urlpath + "index"
	} else {
		urlpath = "/_" + strings.TrimLeft(urlpath, "/")
	}

	hostfilename := path.Join(basePath, path.Clean("./"+protocol+"/"+host+"/"+urlpath+"."+verb+".json"))
	filename := path.Join(basePath, path.Clean("./"+urlpath+"."+verb+".json"))

	if useHostAndProtocol {
		return hostfilename
	} else {
		return filename
	}
}

type FileChecker func(name string) (fi os.FileInfo, err error)

func resolveStatusCode(s string) int {
	prefix := "//! statusCode: "
	if strings.HasPrefix(s, prefix) {
		stringval := strings.TrimLeft(s, prefix)[:3]
		status, err := strconv.Atoi(stringval)
		if err == nil {
			return status
		}
	}
	return 200
}

func stripMetaData(s string) string {
	return s
}
