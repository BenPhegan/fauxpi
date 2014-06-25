package main

import (
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {

	//record := flag.Bool("r", false, "Record request/response pairs to disk.")
	//useHost := flag.Bool("h", false, "Reference files using protocol and host")
	noStubbing := flag.Bool("o", false, "Dont stub any calls")
	port := flag.String("port", "8080", "Port to listen on.")

	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	stubResolver := new(StubResolver)
	fileChecker := os.Stat
	stubResolver.FileChecker = fileChecker

	if !*noStubbing {
		proxy.OnRequest(stubResolver.CheckFilesystemForRequest()).DoFunc(stubResolver.ReturnFileResponse())
	}

	log.Fatal(http.ListenAndServe(":"+*port, proxy))
	fmt.Printf("AND GONE")
}

type ResponseGenerator func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)

func constructFilename(proto string, host string, reqPath string, method string) (string, string) {
	verb := strings.ToLower(method)
	protocol := strings.ToLower(strings.Split(proto, "/")[0])
	urlpath := reqPath
	if strings.HasSuffix(urlpath, "/") {
		urlpath = urlpath + "index"
	} else {
		urlpath = "/_" + strings.TrimLeft(urlpath, "/")
	}

	hostfilename := path.Clean("./" + protocol + "/" + host + "/" + urlpath + "." + verb + ".json")
	filename := path.Clean("./" + urlpath + "." + verb + ".json")

	return hostfilename, filename
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
