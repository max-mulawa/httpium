FROM golang:1.18 as build

RUN mkdir /src
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . ./

RUN mkdir /app 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/httpium ./cmd/server 
#RUN ls /app -la

FROM alpine:latest as deploy
WORKDIR /app
COPY --from=build /app/httpium .
COPY ./config.toml .
COPY ./static/* ./static/

ENTRYPOINT ["/app/httpium"]

EXPOSE 8080