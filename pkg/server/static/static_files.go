package static

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/max-mulawa/httpium/pkg/http/request"
	"github.com/max-mulawa/httpium/pkg/http/response"
	"go.uber.org/zap"
)

type Files struct {
	StaticDir string
	lg        *zap.SugaredLogger
}

func NewStaticFiles(lg *zap.SugaredLogger, staticDir string) *Files {
	return &Files{
		lg:        lg,
		StaticDir: staticDir,
	}
}

func (s *Files) Handle(req *request.HTTPRequest) *response.HTTPResponse {
	filePath := normalize(req.Path)
	filePath = path.Join(s.StaticDir, filePath)

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return response.Response404()
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		s.lg.Error("read file failed", "path", filePath, "err", err)
		return response.Response404()
	}

	res := &response.HTTPResponse{
		Protocol: "HTTP/1.1",
		Code:     response.HTTPCodeOK,
		Headers: map[string]string{
			"Content-Type": "text/html; charset=UTF-8",
		},
		Content: fileContent,
	}

	return res
}

func normalize(reqPath string) string {
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
