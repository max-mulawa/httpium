package request

import (
	"fmt"
	"strings"
)

func Parse(request string) (*Request, error) {

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
			if len(tokens) < 2 {
				continue
			}
			headers[tokens[0]] = tokens[1]
		}
	}
	return headers
}

func parseRequestLine(line string) (*Request, error) {
	tokens := strings.Split(line, " ")
	if len(tokens) < 3 {
		return nil, fmt.Errorf("invalid first request lineh: (%s) with only %d tokens", line, len(tokens))
	}
	r := &Request{}
	r.Method = HttpMethod(tokens[0])
	r.Path = tokens[1]
	r.Protocol = ProtocolVersion(tokens[2])
	return r, nil
}
