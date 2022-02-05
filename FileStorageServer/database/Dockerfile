FROM postgres:12.9-alpine
ENV POSTGRES_USER alex
ENV POSTGRES_DB dbforgolang
ENV POSTGRES_PASSWORD 12345
ADD createTable.sql /docker-entrypoint-initdb.d/

