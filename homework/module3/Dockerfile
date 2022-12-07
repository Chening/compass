FROM golang:1.18 AS build

ENV GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go build -o httpserver

FROM alpine
COPY --from=build /build/httpserver /
EXPOSE 19004

ENTRYPOINT ["/httpserver"]