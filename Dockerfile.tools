FROM docker.io/golang:1.8.1-alpine

RUN apk update
RUN apk add git
RUN apk add wget curl iproute2 tcpdump bridge-utils mtr iperf iftop ldns util-linux ipvsadm ethtool
RUN go get -u github.com/golang/dep/...
