#docker 17.05+
#sudo docker build -t  grpc-demo-golang:1.0 .
#sudo docker run --name grpc-demo-golang -p 8081:50051  --restart always  -d   grpc-demo-golang:1.0 
FROM golang:alpine3.15 as golang
RUN apk update && \
    apk add  bash ca-certificates git && \
    rm -rf /var/cache/apk/*   /tmp/* 
ADD . /data/src/
WORKDIR /data/src/
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn
RUN go mod download
RUN cd ./example/golang/server
RUN go build ./example/golang/server/server.go ./example/golang/server/queue.go


FROM alpine:3.15
RUN TERM=linux && export TERM
USER root
RUN apk update && \
    apk add ca-certificates bash tzdata sudo curl wget openssh git && \ 
    echo "Asia/Shanghai" > /etc/timezone && \
    cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    rm -rf /var/cache/apk/*   /tmp/* 
WORKDIR /data/src
COPY --from=golang /data/src /data/src
CMD ["./server"]
