#!/usr/bin/env bash
#
. ./ENV.sh
docker volume create comments-api-db_postgres

docker run --name comments-api-db \
  -v comments-api-db_postgres:/var/lib/postgresql/data \
  -e POSTGRES_PASSWORD=$DB_PASSWORD -p $DB_PORT:$DB_PORT -d postgres:12.2-alpine 
