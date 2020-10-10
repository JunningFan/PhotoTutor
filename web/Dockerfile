FROM golang:1.13 AS builder

WORKDIR /go/src/app
COPY . .

# China proxy
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.io,direct

# build target 
RUN CGO_ENABLED=0 GOOS=linux  go build  -o main main.go

FROM alpine:latest 
WORKDIR /root/app/ 
COPY --from=builder /go/src/app/main . 

EXPOSE 8080
ENV GIN_MODE release 

# we need a volume to store image
VOLUME ["/root/app/img"]
CMD ["./main"]