FROM ubuntu:18.04



MAINTAINER Alexandr

RUN apt-get -y update && apt-get install -y tzdata

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone



# Обвновление списка пакетов
RUN apt-get -y update

#
# Установка postgresql
#
ENV PGVER 10
RUN apt-get install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER me WITH SUPERUSER PASSWORD 'postgres';" &&\
    createdb -O me forum &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

RUN echo "listen_addresses='*'\nsynchronous_commit = off\nfsync = off\nshared_buffers = 256MB\neffective_cache_size = 1536MB\n" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "wal_buffers = 1MB\nwal_writer_delay = 50ms\nrandom_page_cost = 1.0\nmax_connections = 100\nwork_mem = 8MB\nmaintenance_work_mem = 128MB\ncpu_tuple_cost = 0.0030\ncpu_index_tuple_cost = 0.0010\ncpu_operator_cost = 0.0005" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "full_page_writes = off" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_statement = none" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_duration = off " » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_lock_waits = on" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_min_duration_statement = 5000" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_filename = 'query.log'" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_directory = '/var/log/postgresql'" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_destination = 'csvlog'" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "logging_collector = on" » /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_temp_files = '-1'" » /etc/postgresql/$PGVER/main/postgresql.conf



# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

#
# Сборка проекта
#

RUN apt-get install -y curl
RUN curl —silent —location https://deb.nodesource.com/setup_10.x | bash -
RUN apt-get install -y nodejs
RUN apt-get install -y build-essential




COPY . /forum
WORKDIR /forum

RUN npm install

# Объявлем порт сервера
EXPOSE 5000

# Запускаем, инициализируем базу данных, запускаем приложение
ENV PGPASSWORD postgres
CMD serviceModel postgresql start && psql -h localhost -d forum -U me -p 5432 -a -q -f ./init/init.sql && npm start