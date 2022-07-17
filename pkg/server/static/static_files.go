package static

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/max-mulawa/httpium/pkg/http/request"
	"github.com/max-mulawa/httpium/pkg/http/response"
	"go.uber.org/zap"
)

type StaticFiles struct {
	StaticDir string
	lg        *zap.SugaredLogger
}

func NewStaticFiles(lg *zap.SugaredLogger, staticDir string) *StaticFiles {
	return &StaticFiles{
		lg:        lg,
		StaticDir: staticDir,
	}
}

func (s *StaticFiles) Handle(req *request.HttpRequest) *response.HttpResponse {
	filePath := normalize(req.Path)
	filePath = path.Join(s.StaticDir, filePath)

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return response.Response404()
	}

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		s.lg.Error("read file failed", "path", filePath, "err", err)
		return response.Response404()
	}

	response := &response.HttpResponse{
		Protocol: "HTTP/1.1",
		Code:     200,
		Headers: map[string]string{
			"Content-Type": "text/html; charset=UTF-8",
		},
		Content: fileContent,
	}

	return response
}

func normalize(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimPrefix(path, "\\")

	path = strings.ReplaceAll(path, "../", "")
	path = strings.ReplaceAll(path, "./", "")
	path = strings.ReplaceAll(path, "\\", "/")

	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}

	return path
}
