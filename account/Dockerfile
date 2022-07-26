FROM golang:alpine as builder

WORKDIR /go/src/app

# 在 dev 中获取 reflex 以进行实时重新加载
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
RUN go install github.com/cespare/reflex@latest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./run .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# 从 builder 中复制可执行文件
COPY --from=builder /go/src/app/run .

EXPOSE 8080
CMD ["./run"]