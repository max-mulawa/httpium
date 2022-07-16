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
COPY --from=build /app/httpium .
ENTRYPOINT ["./httpium"]

EXPOSE 8080