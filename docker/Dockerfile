FROM debian

MAINTAINER tengfei-xy

COPY spss_engine /usr/sbin
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
ENTRYPOINT ["spss_engine"]