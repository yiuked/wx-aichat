# 基础镜像
FROM golang:1.18 as builder

WORKDIR /go/src
COPY . .

RUN make install

# 第二个镜像
FROM alpine:latest as ghapi

LABEL MAINTAINER="yiuked"

WORKDIR /go/src

COPY --from=builder /go/src/ghapi ./

ENTRYPOINT ["./ghapi"]
