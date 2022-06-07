FROM golang:1.18.2-alpine
RUN mkdir /src
COPY *.go go.* /src/
COPY pkg /src/pkg
WORKDIR /src
RUN go build -o /notify
ENTRYPOINT ["/notify"]
