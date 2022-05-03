#docker 17.05+
#sudo docker build -t  grpc-demo-golang:1.0 .
#sudo docker run --name grpc-demo-golang -p 8080:50051  --restart always  -d   grpc-demo-golang:1.0 
FROM 1.15.15-alpine3.14 as golang
RUN apk update && \
    apk add  bash ca-certificates git && \
    rm -rf /var/cache/apk/*   /tmp/* 
ADD ./example/golang/server /data/src/
WORKDIR /data/src/
ENV GO111MODULE on
RUN go mod download
RUN go build server.go queue.go


FROM alpine:3.14
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
