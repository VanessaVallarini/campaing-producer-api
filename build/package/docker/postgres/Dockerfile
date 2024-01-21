FROM postgres:14

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=12345
ENV POSTGRES_DB=ads-campaing-db

EXPOSE 5432:5432

ARG LOCAL_SQL=scripts/local
ARG MIGRATION_SQL=scripts/migrations

COPY $LOCAL_SQL /docker-entrypoint-initdb.d/
COPY $MIGRATION_SQL /docker-entrypoint-initdb.d/