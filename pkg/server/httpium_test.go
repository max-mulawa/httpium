package server_test

import (
	"context"
	"fmt"
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

func Test(t *testing.T) {
	filePath := path.Join(tmpContentDir, "index.html")
	err := os.WriteFile(filePath, []byte("<html>ok!</html>"), 0o600)
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Get(fmt.Sprintf("http://localhost:%d/index.html", servicePort))
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
