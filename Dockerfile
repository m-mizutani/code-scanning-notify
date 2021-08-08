FROM golang
RUN go build -o notify
CMD ["./notify"]
