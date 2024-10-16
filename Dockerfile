FROM alpine

LABEL MAINTAINER tengfei-xy

RUN mkdir /data
WORKDIR /data

COPY spss_engine /usr/sbin
ENTRYPOINT ["spss_engine"]