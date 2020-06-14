FROM ubuntu:18.04

MAINTAINER Alexandra Kruchinina

# Обвновление списка пакетов
RUN apt-get -y update

# Установка postgresql
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

# config Postgre
RUN echo "synchronous_commit = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = off" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

#
# Сборка проекта
#

# Установка Nodejs & npm
# install from nodesource using apt-get
# https://www.digitalocean.com/community/tutorials/how-to-install-node-js-on-an-ubuntu-14-04-server
RUN apt-get update && apt-get install -y \
curl
RUN curl -sL https://deb.nodesource.com/setup_9.x | bash -
RUN apt-get install -y nodejs


# Копируем исходный код в Docker-контейнер
ADD . /db
WORKDIR /db

# Собираем и устанавливаем пакет
RUN npm install

# Объявлем порт сервера
EXPOSE 5000

#
# Запускаем PostgreSQL и сервер


ENV PGPASSWORD postgres
CMD service postgresql start && psql -h localhost -d forum -U me -p 5432 -a -q -f ./init/init.sql && npm startg