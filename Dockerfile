FROM golang:1.16.5-buster AS build
COPY . /opt/server
ENV GO111MODULE=on
WORKDIR /opt/server/
RUN go mod tidy
RUN go build
FROM ubuntu:20.04 AS release
RUN apt -y update && apt install -y locales gnupg2
RUN locale-gen en_US.UTF-8
RUN update-locale LANG=en_US.UTF-8
ENV PGVER 12
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y && apt-get install -y postgresql postgresql-contrib
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER dbmsmaster WITH PASSWORD 'dbms' SUPERUSER;" &&\
    createdb -O dbmsmaster dbmsforum &&\
    /etc/init.d/postgresql stop
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
USER root
EXPOSE 5000
WORKDIR /usr/bin
COPY --from=build /opt/server/technopark-dbms /usr/bin/
COPY ./init.sql /usr/bin/
ENV PGPASSWORD dbms
CMD service postgresql start && psql --quiet -h localhost -p 5432 -U dbmsmaster -d dbmsforum -a -f init.sql && technopark-dbms
