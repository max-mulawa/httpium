package static

import (
	"errors"
	"os"
	"strings"

	"github.com/max-mulawa/httpium/pkg/http/request"
	"github.com/max-mulawa/httpium/pkg/http/response"
	"go.uber.org/zap"
)

type Files struct {
	StaticDir    string
	DefaultFiles []string
	lg           *zap.SugaredLogger
}

func NewStaticFiles(lg *zap.SugaredLogger, staticDir string, defaultFiles []string) *Files {
	return &Files{
		lg:           lg,
		StaticDir:    staticDir,
		DefaultFiles: defaultFiles,
	}
}

func (s *Files) Handle(req *request.HTTPRequest) *response.HTTPResponse {
	relPath := normalize(req.Path)

	filePath, err := s.getLocalPath(relPath)
	if errors.Is(err, ErrFileNotFound) {
		return response.Response404()
	}

	if err != nil {
		s.lg.Error("get local file path failed", "path", relPath, "err", err)
		return response.Response500()
	}

	if filePath == "" {
		return okResponse([]byte("<html><body>Welcome to the world of httpium</body></html>"))
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		s.lg.Error("read file failed", "path", filePath, "err", err)
		return response.Response500()
	}

	res := okResponse(fileContent)

	return res
}

func okResponse(content []byte) *response.HTTPResponse {
	return &response.HTTPResponse{
		Protocol: "HTTP/1.1",
		Code:     response.HTTPCodeOK,
		Headers: map[string]string{
			"Content-Type": "text/html; charset=UTF-8",
		},
		Content: content,
	}
}

func normalize(reqPath string) string {
	if isDefaultPath(reqPath) {
		return reqPath
	}

	reqPath = strings.TrimPrefix(reqPath, "/")
	reqPath = strings.TrimPrefix(reqPath, "\\")

	reqPath = strings.ReplaceAll(reqPath, "../", "")
	reqPath = strings.ReplaceAll(reqPath, "./", "")
	reqPath = strings.ReplaceAll(reqPath, "\\", "/")

	if idx := strings.Index(reqPath, "?"); idx != -1 {
		reqPath = reqPath[:idx]
	}

	return reqPath
}
