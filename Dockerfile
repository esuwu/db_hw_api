FROM ubuntu:20.04

MAINTAINER Alexandr Dolgavin

RUN apt-get -y update && apt-get install -y tzdata

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone


ENV PGVER 12
RUN apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER me WITH SUPERUSER PASSWORD 'postgres';" &&\
    createdb -O me forum &&\
    /etc/init.d/postgresql stop

RUN echo "synchronous_commit = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = off" >> /etc/postgresql/$PGVER/main/postgresql.conf

RUN echo "listen_addresses='*'\nsynchronous_commit = off\nfsync = off\nshared_buffers = 256MB\neffective_cache_size = 1536MB\n" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "wal_buffers = 1MB\nwal_writer_delay = 50ms\nrandom_page_cost = 1.0\nmax_connections = 100\nwork_mem = 8MB\nmaintenance_work_mem = 128MB\ncpu_tuple_cost = 0.0030\ncpu_index_tuple_cost = 0.0010\ncpu_operator_cost = 0.0005" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "full_page_writes = off" >>  /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_statement = none" >>  /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_duration = off " >>  /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_lock_waits = on" >>  /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_min_duration_statement = 5000" >>  /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_filename = 'query.log'" >>  /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_directory = '/var/log/postgresql'" >>  /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_destination = 'csvlog'" >>  /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "logging_collector = on" >>  /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_temp_files = '-1'" >>  /etc/postgresql/$PGVER/main/postgresql.conf

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

#
# Сборка проекта
#

RUN apt-get update && apt-get install -y \
curl
RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get install -y nodejs



ADD . /db
WORKDIR /db

# Собираем и устанавливаем пакет
RUN npm install

# Объявлем порт сервера
EXPOSE 5000

#
# Запускаем PostgreSQL и сервер


ENV PGPASSWORD postgres
CMD service postgresql start && psql -h localhost -d forum -U me -p 5432 -a -q -f ./init/init.sql && npm start