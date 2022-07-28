package request

import (
	"fmt"
	"strings"

	"github.com/max-mulawa/httpium/pkg/http"
)

var (
	requestLineTokensCnt   int = 3
	requestHeaderTokensCnt int = 2
)

func Parse(request string) (*HTTPRequest, error) {
	lines := strings.Split(request, "\r\n")
	reqLineField := lines[0]
	r, err := parseRequestLine(reqLineField)

	if err != nil {
		return nil, err
	}

	headers := parseHeadersLines(lines[1:])

	r.Headers = headers

	return r, nil
}

func parseHeadersLines(headerLiens []string) map[string]string {
	headers := make(map[string]string, len(headerLiens))

	for _, line := range headerLiens {
		if strings.Contains(line, ": ") {
			tokens := strings.Split(line, ": ")
			if len(tokens) < requestHeaderTokensCnt {
				continue
			}

			headers[tokens[0]] = tokens[1]
		}
	}

	return headers
}

func parseRequestLine(line string) (*HTTPRequest, error) {
	tokens := strings.Split(line, " ")

	if len(tokens) < requestLineTokensCnt {
		return nil, fmt.Errorf("invalid first request lineh: (%s) with only %d tokens", line, len(tokens))
	}

	r := &HTTPRequest{
		Method:   HTTPMethod(tokens[0]),
		Path:     tokens[1],
		Protocol: http.ProtocolVersion(tokens[2]),
	}

	return r, nil
}
