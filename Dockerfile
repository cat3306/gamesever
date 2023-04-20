FROM golang:1.18-buster AS build
WORKDIR /app
COPY ./ ./
RUN go env -w GOPROXY=https://goproxy.cn
RUN go mod download
RUN go build -o /gameserver

FROM  centos AS build-release-stage
ENV TZ Asia/Shanghai
COPY --from=build /gameserver /gameserver
COPY conf/conf.json /conf.json
COPY private_key.pem /private_key.pem
WORKDIR /
#RUN chmod u+x /gameserver &&  mkdir -p /var/log/go_log

CMD ["/gameserver","-c","conf.json"]
