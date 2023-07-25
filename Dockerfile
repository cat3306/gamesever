FROM ubuntu
ENV TZ Asia/Shanghai
COPY gameserver /gameserver
COPY conf/conf.json /conf.json
COPY private_key.pem /private_key.pem
WORKDIR /
#RUN chmod u+x /gameserver &&  mkdir -p /var/log/go_log

CMD ["/gameserver","-c","conf.json"]
