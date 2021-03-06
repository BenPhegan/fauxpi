package main

import (
	"errors"
	"github.com/elazarl/goproxy"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
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
	slash := string(os.PathSeparator)
	verb := strings.ToLower(method)
	protocol := strings.ToLower(strings.Split(proto, "/")[0])
	urlpath := reqPath
	if strings.HasSuffix(urlpath, "/") {
		urlpath = urlpath + "index"
	} else {
		urlpath = "/_" + strings.TrimLeft(urlpath, "/")
	}

	hostfilename := path.Join(basePath, path.Clean("."+slash+protocol+slash+host+slash+urlpath+"."+verb+".json"))
	filename := path.Join(basePath, path.Clean("."+slash+urlpath+"."+verb+".json"))

	if useHostAndProtocol {
		return hostfilename
	} else {
		return filename
	}
}

type FileChecker func(name string) (fi os.FileInfo, err error)

func resolveStatusCode(s string) int {
	regex := regexp.MustCompile(`//!.+statusCode:\s*(\d{3})`)
	result := regex.FindStringSubmatch(s)
	if len(result) > 0 {
		intresult, _ := strconv.Atoi(result[1])
		return intresult
	}
	return 200
}

func stripMetaData(s string) string {
	if strings.HasPrefix(s, "//!") {
		stringarray := strings.SplitAfterN(s, "\n", 2)
		return stringarray[1]
	}
	return s
}

type ResponseFunc func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response

func (sr StubResolver) RecordResponse() ResponseFunc {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		var filename = constructFilename(ctx.Req.Proto, ctx.Req.Host, ctx.Req.URL.Path, ctx.Req.Method, sr.UseHostAndProtocol, sr.StubRoot)
		sr.saveResponse(resp, ctx, filename)
		return resp
	}
}

func (sr StubResolver) saveResponse(resp *http.Response, ctx *goproxy.ProxyCtx, filename string) {
	if resp != nil {
		//TODO Cleanup error
		ctx.Logf("Writing response to here: " + filename)
		if _, err := sr.FileChecker(filepath.Dir(filename)); err != nil {
			os.MkdirAll(filepath.Dir(filename), 0777)
		}
		f, _ := os.Create(filename)
		defer f.Close()

		resp.Body = NewTeeReadCloser(resp.Body, NewFileStream(filename))
		f.Sync()
	}
}

type TeeReadCloser struct {
	r io.Reader
	w io.WriteCloser
	c io.Closer
}

func NewTeeReadCloser(r io.ReadCloser, w io.WriteCloser) io.ReadCloser {
	return &TeeReadCloser{io.TeeReader(r, w), w, r}
}

func (t *TeeReadCloser) Read(b []byte) (int, error) {
	return t.r.Read(b)
}

func (t *TeeReadCloser) Close() error {
	err1 := t.c.Close()
	err2 := t.w.Close()
	if err1 == nil && err2 == nil {
		return nil
	}
	if err1 != nil {
		return err2
	}
	return err1
}

type FileStream struct {
	path string
	f    *os.File
}

func NewFileStream(path string) *FileStream {
	return &FileStream{path, nil}
}

func (fs *FileStream) Write(b []byte) (nr int, err error) {
	if fs.f == nil {
		fs.f, err = os.Create(fs.path)
		if err != nil {
			return 0, err
		}
	}
	return fs.f.Write(b)
}

func (fs *FileStream) Close() error {
	if fs.f == nil {
		return errors.New("FileStream was never written into")
	}
	return fs.f.Close()
}
