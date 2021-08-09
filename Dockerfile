FROM golang:1.16-alpine
RUN mkdir /src
COPY *.go go.* /src/
COPY pkg /src/pkg
WORKDIR /src
RUN go build -o /notify
CMD ["/notify"]
