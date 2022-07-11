package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRequestLine(t *testing.T) {
	requestLine := "GET /hello.txt HTTP/1.1"
	request, err := parseRequestLine(requestLine)

	require.NoError(t, err)
	require.Equal(t, "GET", string(request.Method))
	require.Equal(t, "HTTP/1.1", string(request.Protocol))
	require.Equal(t, "/hello.txt", request.Path)
}

func TestParseHeadersLines(t *testing.T) {
	headerLines := []string{
		"User-Agent: curl/7.16.3 libcurl/7.16.3 OpenSSL/0.9.7l zlib/1.2.3",
		"Host: www.example.com",
		"Accept-Language: en, mi",
	}

	headers := parseHeadersLines(headerLines)
	require.Contains(t, headers, "User-Agent")
	require.Contains(t, headers, "Host")
	require.Contains(t, headers, "Accept-Language")
}
