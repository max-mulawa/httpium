# httpium

httpium is a simple http server written in go.

[![Go](https://github.com/max-mulawa/httpium/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/max-mulawa/httpium/actions/workflows/go.yml)

# get started

This will bring up httpium server listening on port 8080

Run locally
```bash
make build
./bin/httpium
```

Run on docker
```bash
make docker-run
```

Run on Kubernetes
```bash
kubectl run --image docker.io/mulawam/httpium:latest httpium
kubectl port-forward po/httpium 8080:8080

curl localhost:8080
```

# debug

Send SIGTERM to debugging code in vscode
```bash
pgrep debug | xargs kill -s 15
```

Trigger daemon reload (triggers ExecReload that in return sends SIGHUP signal to the process)
```bash
systemctl reload httpium.service
```

Connect to running `httpium` container
```bash
# run 
make docker-run

# attach
docker exec -it httpium sh
```

