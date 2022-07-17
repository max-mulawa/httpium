package request

import "github.com/max-mulawa/httpium/pkg/http"

type HttpRequest struct {
	Method    HttpMethod
	UserAgent string
	Path      string
	Protocol  http.ProtocolVersion
	Headers   map[string]string
}

type HttpMethod string
