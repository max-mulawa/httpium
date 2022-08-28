package server_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/max-mulawa/httpium/pkg/server"
	"github.com/max-mulawa/httpium/pkg/server/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"net/http"
)

var (
	tmpContentDir   string
	servicePort     uint = 8090
	client               = http.Client{}
	contentFileName      = "index.html"
	defaultFileName      = "default.html"
)

func TestMain(m *testing.M) {
	currDir, _ := os.Getwd()
	tmpContentDir, _ = ioutil.TempDir(currDir, "httpium-content")

	defer os.RemoveAll(tmpContentDir)

	cfg := &config.HttpiumConfig{
		Server: config.ServerOptions{
			Port: servicePort,
		},
		Content: config.ContentOptions{
			StaticDir: tmpContentDir,
			Default:   []string{defaultFileName},
		},
	}
	srv := server.NewServer(context.Background(), zap.NewNop().Sugar(), cfg)

	var err error

	go func() { err = srv.Start() }()

	if err != nil {
		log.Printf("error on start: %v", err)
		return
	}

	m.Run()

	err = srv.Stop()
	if err != nil {
		log.Printf("error on stop: %v", err)
		return
	}
}

func TestStaticContent(t *testing.T) {
	t.Parallel()
	t.Run("existing file serverd", func(t *testing.T) {
		fileContent := "<html>ok!</html>"
		writeTmpFile(t, contentFileName, fileContent)

		resp, err := client.Get(fmt.Sprintf("http://localhost:%d/%s", servicePort, contentFileName))
		require.NoError(t, err)

		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		responseContentEquals(t, resp, fileContent)
	})

	t.Run("default file served", func(t *testing.T) {
		fileContent := "<html>default!</html>"
		writeTmpFile(t, defaultFileName, fileContent)

		resp, err := client.Get(fmt.Sprintf("http://localhost:%d/", servicePort))
		require.NoError(t, err)

		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		responseContentEquals(t, resp, fileContent)
	})

	t.Run("for missing file 404 error code returned", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("http://localhost:%d/index2.html", servicePort))
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func responseContentEquals(t *testing.T, resp *http.Response, fileContent string) {
	t.Helper()

	content, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, fileContent, string(content))
}

func writeTmpFile(t *testing.T, fileName, fileContent string) {
	t.Helper()

	filePath := path.Join(tmpContentDir, fileName)
	err := os.WriteFile(filePath, []byte(fileContent), 0o600)
	require.NoError(t, err)
}
