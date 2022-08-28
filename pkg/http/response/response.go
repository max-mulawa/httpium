package response

import (
	"fmt"
	"strings"

	"github.com/max-mulawa/httpium/pkg/http"
)

type HTTPResponse struct {
	Protocol http.ProtocolVersion
	Headers  map[string]string
	Content  []byte
	Code     HTTPResponseCode
}

type HTTPResponseCode uint

var (
	HTTPCodeOK                HTTPResponseCode = 200
	HTTPCodeNotFound          HTTPResponseCode = 404
	HTTPCodeInternalServerErr HTTPResponseCode = 500
)

var codeAsText map[HTTPResponseCode]string

func init() {
	codeAsText = map[HTTPResponseCode]string{
		HTTPResponseCode(uint(HTTPCodeOK)):                "OK",
		HTTPResponseCode(uint(HTTPCodeNotFound)):          "Not Found",
		HTTPResponseCode(uint(HTTPCodeInternalServerErr)): "Internal Server Error",
	}
}

func (r *HTTPResponse) Build() ([]byte, error) {
	builder := strings.Builder{}

	_, err := builder.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Protocol, r.Code, r.Code.getCodeAsText()))
	if err != nil {
		return nil, err
	}

	for hname, hval := range r.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", hname, hval))
	}

	_, err = builder.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(r.Content)))
	if err != nil {
		return nil, err
	}

	if len(r.Content) > 0 {
		_, err = builder.WriteString("\r\n")
		if err != nil {
			return nil, err
		}

		_, err = builder.Write(r.Content)
		if err != nil {
			return nil, err
		}
	}

	return []byte(builder.String()), nil
}

func Response404() *HTTPResponse {
	return &HTTPResponse{
		Protocol: "HTTP/1.1",
		Code:     HTTPCodeNotFound,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Content: []byte("<html>Page not found</html>"),
	}
}

func Response500() *HTTPResponse {
	return &HTTPResponse{
		Protocol: "HTTP/1.1",
		Code:     HTTPCodeInternalServerErr,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Content: []byte("<html>Internal Server Error occurred</html>"),
	}
}

func (c HTTPResponseCode) getCodeAsText() string {
	if v, ok := codeAsText[c]; ok {
		return v
	}

	return ""
}
