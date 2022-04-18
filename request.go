package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/armosec/utils-go/httputils"
)

type FlagParser struct {
	fullURL       url.URL
	pathToBody    string
	body          string
	pathToHeaders string
	headers       string
	method        string
}

func NewFlagParser() *FlagParser {
	return &FlagParser{
		fullURL: url.URL{},
	}
}
func (f *FlagParser) parser() {

	flag.StringVar(&f.method, "method", "", "http method (GET/POST/DELETE)")
	flag.StringVar(&f.fullURL.Scheme, "scheme", "http", "request scheme")
	flag.StringVar(&f.fullURL.Host, "host", "", "host")
	flag.StringVar(&f.fullURL.Path, "path", "", "path")

	// flag.StringVar(&f.pathToBody, "body", "", "body")
	flag.StringVar(&f.pathToBody, "path-body", "", "path to body")
	flag.StringVar(&f.pathToHeaders, "headers", "", "http headers")
	// flag.StringVar(&f.pathToHeaders, "path-headers", "", "path to headers")

	flag.Parse()
}

func (f *FlagParser) validate() error {
	if f.fullURL.Host == "" {
		return fmt.Errorf("missing host")
	}
	if f.method == "" {
		return fmt.Errorf("missing method")
	}
	return nil
}

func loadHeaders(f *FlagParser) (map[string]string, error) {
	headers := map[string]string{}
	if f.headers != "" {
		splitteedHeaders := strings.Split(f.headers, ";")
		for i := range splitteedHeaders {
			header := strings.Split(splitteedHeaders[i], ":")
			if len(header) == 2 {
				headers[header[0]] = strings.TrimLeft(header[1], " ")
			}
		}
		return headers, nil
	}
	if f.pathToHeaders != "" {
		// Not supported
	}
	return headers, nil
}

func loadBody(f *FlagParser) ([]byte, error) {
	if f.body != "" {
		// Not suppored
	}
	if f.pathToBody != "" {
		fmt.Printf("loading body from: %s\n", f.pathToBody)
		return os.ReadFile(f.pathToBody)
	}
	return []byte{}, nil
}

// Request run a http request
func Request(f *FlagParser) error {
	var resp *http.Response
	var err error

	headers, e := loadHeaders(f)
	if e != nil {
		return e
	}
	body, e := loadBody(f)
	if e != nil {
		return e
	}

	fmt.Printf("method: %s, url: %s, headers: %v, body: %s\n", f.method, f.fullURL.String(), headers, body)

	switch f.method {
	case "POST", "post":
		resp, err = httputils.HttpPost(http.DefaultClient, f.fullURL.String(), headers, body)
	case "GET", "get":
		resp, err = httputils.HttpGet(http.DefaultClient, f.fullURL.String(), headers)
	case "DELETE", "delete":
		resp, err = httputils.HttpGet(http.DefaultClient, f.fullURL.String(), headers)
	default:
		return fmt.Errorf("method %s not supported", f.method)
	}

	if err != nil {
		return err
	}

	strResp, e := httputils.HttpRespToString(resp)
	if e != nil {
		return fmt.Errorf("failed to parse http response to string, reason: %s", e.Error())
	} else {
		fmt.Printf("response: %s", strResp)
	}
	return nil
}
