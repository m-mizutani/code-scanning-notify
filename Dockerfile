FROM golang:1.16-alpine
RUN mkdir /app
WORKDIR /app
COPY *.go ./
COPY go.* ./
RUN go build -o /notify
CMD ["/notify"]
