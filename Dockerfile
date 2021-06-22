FROM golang:1.16.5-buster AS build

# Копируем исходный код в Docker-контейнер
COPY . /opt/server

# Собираем генераторы
ENV GO111MODULE=on

WORKDIR /opt/server/
# RUN go mod tidy
RUN go build

FROM ubuntu:20.04 AS release

# Make the "en_US.UTF-8" locale so postgres will be utf-8 enabled by default
RUN apt -y update && apt install -y locales gnupg2
RUN locale-gen en_US.UTF-8
RUN update-locale LANG=en_US.UTF-8

#
# Install postgresql
#

ENV PGVER 12
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y && apt-get install -y postgresql postgresql-contrib

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER dbmsmaster WITH PASSWORD 'dbms' SUPERUSER;" &&\
    createdb -O dbmsmaster dbmsforum &&\
    /etc/init.d/postgresql stop
#
## Adjust PostgreSQL configuration so that remote connections to the
## database are possible.
#RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
#
## And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
#RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
#
#EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# Объявлем порт сервера
EXPOSE 5000

WORKDIR /usr/bin

# Собранный ранее сервер
COPY --from=build /opt/server/technopark-dbms /usr/bin/
COPY ./init.sql /usr/bin/

ENV PGPASSWORD dbms
CMD service postgresql start && psql --quiet -h localhost -p 5432 -U dbmsmaster -d dbmsforum -a -f init.sql && technopark-dbms