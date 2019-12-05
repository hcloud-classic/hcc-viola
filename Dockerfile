# Dockerfile
FROM ubuntu:latest
MAINTAINER ish <ish@innogrid.com>

RUN mkdir -p /cello/
WORKDIR /cello/

ADD flute /cello/
RUN chmod 755 /cello/cello

EXPOSE 7000

CMD ["/cello/cello"]
