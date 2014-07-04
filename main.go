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

	record := flag.Bool("r", false, "Record request/response pairs to disk.")
	useHost := flag.Bool("h", false, "Reference files using protocol and host")
	noStubbing := flag.Bool("o", false, "Dont stub any calls")
	port := flag.String("port", "8080", "Port to listen on.")
	stubRoot := flag.String("d", "", "The directory root for your stubs")

	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	//Get the root directory for stub files
	if *stubRoot == "" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		stubRoot = &pwd
	}

	stubResolver := new(StubResolver)
	fileChecker := os.Stat
	stubResolver.FileChecker = fileChecker
	stubResolver.UseHostAndProtocol = *useHost
	stubResolver.StubRoot = *stubRoot

	if !*noStubbing {
		proxy.OnRequest(stubResolver.CheckFilesystemForRequest()).DoFunc(stubResolver.ReturnFileResponse())
	}

	if *record {
		proxy.OnResponse().DoFunc(stubResolver.RecordResponse())
	}

	log.Fatal(http.ListenAndServe(":"+*port, proxy))
	fmt.Printf("AND GONE")
}
