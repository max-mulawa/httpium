package request

type Request struct {
	Method    HttpMethod
	UserAgent string
	Path      string
	Protocol  ProtocolVersion
	Headers   map[string]string
}

type HttpMethod string

type ProtocolVersion string
