package main

import (
	"net/url"
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
