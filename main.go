package main

import (
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"os"
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
