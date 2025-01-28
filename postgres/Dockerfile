FROM postgres:16-alpine
COPY sqlite3_dump_edited.sql /docker-entrypoint-initdb.d/

ENV POSTGRES_USER=postgres
#this should be RUN --mount=type=secret,POSTGRES_PASSWORD=ARG_PG_PASS
ENV POSTGRES_PASSWORD=postgres
ENV PGDATA=/data

#RUN ["/usr/local/bin/docker-entrypoint.sh", "postgres"]


