package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/armosec/utils-go/httputils"
)

type FlagParser struct {
	fullURL      url.URL
	pathToBody   string
	headers      string
	method       string
	pathToOutput string
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

	flag.StringVar(&f.pathToBody, "path-body", "", "path to body")
	flag.StringVar(&f.headers, "headers", "", "http headers")
	flag.StringVar(&f.pathToOutput, "path-output", "", "path to output file")

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
	return headers, nil
}

func loadBody(f *FlagParser) ([]byte, error) {
	if f.pathToBody != "" {
		fmt.Printf("loading body from: %s\n", f.pathToBody)
		return os.ReadFile(f.pathToBody)
	}
	return []byte{}, nil
}

func setHeaders(req *http.Request, headers map[string]string) {
	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
}

// Request run a http request
func Request(f *FlagParser) (string, error) {
	var resp *http.Response
	var err error

	headers, e := loadHeaders(f)
	if e != nil {
		return "", e
	}
	body, e := loadBody(f)
	if e != nil {
		return "", e
	}

	fmt.Printf("method: %s, url: %s, headers: %v, body: %s\n", f.method, f.fullURL.String(), headers, body)

	methods := []string{http.MethodPost, http.MethodGet, http.MethodDelete}
	if slices.Contains(methods, strings.ToUpper(f.method)) {
		req, err := http.NewRequest(f.method, f.fullURL.String(), bytes.NewReader(body))
		if err != nil {
			return "", err
		}
		setHeaders(req, headers)
		resp, err = http.DefaultClient.Do(req)
	} else {
		return "", fmt.Errorf("method %s not supported", f.method)
	}

	if err != nil {
		return "", err
	}

	strResp, e := httputils.HttpRespToString(resp)
	if e != nil {
		return "", fmt.Errorf("failed to parse http response to string, reason: %s", e.Error())
	}

	if f.pathToOutput != "" {
		if err := os.WriteFile(f.pathToOutput, []byte(strResp), 0644); err != nil {
			return "", fmt.Errorf("error writing response to file: %s", err.Error())
		}
		fmt.Printf("response was written to file: %s\n", f.pathToOutput)

	}
	return strResp, nil
}
