FROM centos
ENV TZ Asia/Shanghai
LABEL name="gameserver" author="cat3306" branch="master"
ADD  gameserver /gameserver
ADD conf/conf.json /conf.json
WORKDIR /
RUN chmod u+x /gameserver &&  mkdir -p /var/log/go_log
CMD ["/gameserver","-c","conf.json"]