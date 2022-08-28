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
	tmpContentDir string
	servicePort   uint = 8090
	client             = http.Client{}
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
	t.Run("existing file serverd", func(t *testing.T) {
		filePath := path.Join(tmpContentDir, "index.html")
		fileContent := "<html>ok!</html>"
		err := os.WriteFile(filePath, []byte(fileContent), 0o600)
		require.NoError(t, err)

		resp, err := client.Get(fmt.Sprintf("http://localhost:%d/index.html", servicePort))
		require.NoError(t, err)

		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		content, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, fileContent, string(content))
	})

	t.Run("for missing file 404 error code returned", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("http://localhost:%d/index2.html", servicePort))
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
