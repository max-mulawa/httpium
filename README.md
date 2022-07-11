Build docker
```
docker build . -t localhost/httpium:dev
```

Run httpium on port 8080
```
docker run --rm --name httpium -p 8080:8080 localhost/httpium:dev
```