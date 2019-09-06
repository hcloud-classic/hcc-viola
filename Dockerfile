# Dockerfile
FROM ubuntu:latest
MAINTAINER ish <ish@innogrid.com>

RUN mkdir -p /GraphQL_Cello/
WORKDIR /GraphQL_Cello/

ADD GraphQL_Cello /GraphQL_Cello/
RUN chmod 755 /GraphQL_Cello/GraphQL_Cello

EXPOSE 8001

CMD ["/GraphQL_Cello/GraphQL_Cello"]
