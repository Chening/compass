FROM golang:1.18 AS build

WORKDIR /build
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o httpserver main.go

FROM alpine
COPY --from=build /build/httpserver /
EXPOSE 19004

ENTRYPOINT ["/httpserver"]