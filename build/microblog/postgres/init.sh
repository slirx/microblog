#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE microblog_auth;
    CREATE DATABASE microblog_post;
    CREATE DATABASE microblog_registration;
    CREATE DATABASE microblog_user;
EOSQL
