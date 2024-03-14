package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	{
		f := FlagParser{
			fullURL: url.URL{
				Host: "localhost:8080",
			},
		}
		f.method = "POST"
		assert.NoError(t, f.validate())
	}
	{
		f := FlagParser{
			fullURL: url.URL{
				Host: "",
			},
		}
		assert.Error(t, f.validate())
	}
}
func TestLoadHeaders(t *testing.T) {
	{
		f := FlagParser{
			headers: "Content-Type:application/json",
		}
		headers, err := loadHeaders(&f)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(headers))
		assert.Equal(t, "application/json", headers["Content-Type"])

	}
	{
		f := FlagParser{
			headers: "Content-Type:application/json;Auth: bar",
		}
		headers, err := loadHeaders(&f)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(headers))
		assert.Equal(t, "application/json", headers["Content-Type"])
		assert.Equal(t, "bar", headers["Auth"])

	}
}

func setupTestServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("GET method received"))
		case http.MethodPost:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("POST method received"))
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("DELETE method received"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("Unsupported method"))
		}
	}))
	t.Cleanup(func() { server.Close() })
	return server
}

func TestRequest(t *testing.T) {
	server := setupTestServer(t)

	testcases := map[string]struct {
		method     string
		response   string
		outputFile string
	}{
		"GET": {
			method:     "get",
			response:   "GET method received",
			outputFile: "test_output.txt",
		},
		"POST": {
			method:   http.MethodPost,
			response: "POST method received",
		},
		"DELETE": {
			method:   http.MethodDelete,
			response: "DELETE method received",
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			flagParser := NewFlagParser()
			flagParser.method = tc.method
			flagParser.fullURL.Host = server.URL[7:]
			flagParser.fullURL.Scheme = "http"
			flagParser.fullURL.Path = "/"
			flagParser.headers = "Header1:Value1;Header2:Value2"
			if tc.outputFile != "" {
				flagParser.pathToOutput = tc.outputFile
				t.Cleanup(func() { _ = os.Remove(flagParser.pathToOutput) })
			}
			strResp, err := Request(flagParser)
			assert.NoErrorf(t, err, "request failed")
			assert.Equalf(t, tc.response, strResp, "response mismatch")

			if tc.outputFile != "" {
				data, err := os.ReadFile(flagParser.pathToOutput)
				assert.NoErrorf(t, err, "failed to read output file")
				assert.Equalf(t, tc.response, string(data), "response mismatch in output file")
			}
		})
	}

}
