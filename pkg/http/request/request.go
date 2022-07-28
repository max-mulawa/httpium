package request

import "github.com/max-mulawa/httpium/pkg/http"

type HTTPRequest struct {
	Method    HTTPMethod
	UserAgent string
	Path      string
	Protocol  http.ProtocolVersion
	Headers   map[string]string
}

type HTTPMethod string
