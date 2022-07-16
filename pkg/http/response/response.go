package response

import (
	"fmt"
	"strings"

	"github.com/max-mulawa/httpium/pkg/http"
)

type HttpResponse struct {
	Protocol http.ProtocolVersion
	Headers  map[string]string
	Content  []byte
	Code     HttpResponseCode
}

type HttpResponseCode uint

var codeAsText map[HttpResponseCode]string

func init() {
	codeAsText = map[HttpResponseCode]string{
		HttpResponseCode(200): "OK",
	}
}

func (r *HttpResponse) Build() ([]byte, error) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Protocol, r.Code, r.Code.getCodeAsText()))
	for hname, hval := range r.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", hname, hval))
	}
	builder.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(r.Content)))
	if len(r.Content) > 0 {
		builder.WriteString("\r\n")
		builder.Write(r.Content)
	}
	//"HTTP/1.1 200 OK\r\nContent-Length: 12\r\nContent-Type: text/plain\r\n\r\nHello World!"
	return []byte(builder.String()), nil
}

func (c HttpResponseCode) getCodeAsText() string {
	if v, ok := codeAsText[c]; ok {
		return v
	}
	return ""
}
