package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	stubResolver := new(StubResolver)
	stubSoure := new(StubSource)
	fileChecker := os.Stat
	stubResolver.StubSource = stubSoure
	stubResolver.FileChecker = fileChecker
	proxy.OnRequest(stubResolver.CheckFilesystemForRequest()).DoFunc(stubResolver.ReturnResponse())

	log.Fatal(http.ListenAndServe(":8090", proxy))
	fmt.Printf("AND GONE")
}

type StubResolver struct {
	StubSource  *StubSource
	FileChecker FileChecker
}

type StubSource struct {
}

func (ss StubSource) Stubs() []Stub {
	var stubs []Stub
	stubs = append(stubs, Stub{MatchType: LiteralMatch, Match: "http://google.com/", Body: "GOOGLE"})
	return stubs
}

func (sr StubResolver) CheckUrlAgainstStub() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		for _, s := range sr.StubSource.Stubs() {
			if s.Match == req.RequestURI {
				return true
			}
		}
		return false
	}
}

type ResponseGenerator func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)

func (sr StubResolver) CheckFilesystemForRequest() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		ctx.Logf(req.Proto + "," + req.Host + "," + req.URL.Path + "," + req.Method)
		var hostfilename, filename = constructFilename(req.Proto, req.Host, req.URL.Path, req.Method)

		ctx.Logf(hostfilename)
		if _, err := sr.FileChecker(filename); err == nil {
			ctx.Logf("***FOUND FILE***")
			return true
		}
		return false
	}
}

func constructFilename(proto string, host string, reqPath string, method string) (string, string) {
	verb := strings.ToLower(method)
	protocol := strings.ToLower(strings.Split(proto, "/")[0])
	urlpath := reqPath
	if strings.HasSuffix(urlpath, "/") {
		urlpath = urlpath + "index"
	} else {
		//parts := strings.Split(urlpath, "/")
		//newparts := len(parts)
		urlpath = "/_" + strings.TrimLeft(urlpath, "/")
	}

	hostfilename := path.Clean("./" + protocol + "/" + host + "/" + urlpath + "." + verb + ".json")
	filename := path.Clean("./" + urlpath + "." + verb + ".json")

	return hostfilename, filename
}

type FileChecker func(name string) (fi os.FileInfo, err error)

func (sr StubResolver) ReturnResponse() ResponseGenerator {
	return func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		for _, s := range sr.StubSource.Stubs() {
			return r, goproxy.NewResponse(r,
				goproxy.ContentTypeText, http.StatusForbidden,
				s.Body)
			return r, nil
		}

		return nil, nil
	}
}

type Stub struct {
	MatchType MatchType
	Match     string
	Headers   []string
	Body      string
}

type MatchType int

const (
	LiteralMatch MatchType = iota
	RegexMatch
)
