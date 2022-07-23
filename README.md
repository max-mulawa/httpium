# httpium

httpium is a simple http server written in go.

[![Go](https://github.com/max-mulawa/httpium/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/max-mulawa/httpium/actions/workflows/go.yml)

# get started

This will bring up httpium server listening on port 8080
```console
./bin/httpium
```

# debug

Send SIGTERM to debugging code in vscode
```console
pgrep debug | xargs kill -s 15
```

Trigger daemon reload (triggers ExecReload that in return sends SIGHUP signal to the process)
```console
systemctl reload httpium.service
```
