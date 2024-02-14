FROM postgres:16

COPY *.sql /docker-entrypoint-initdb.d/

ENV POSTGRES_DB=main
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=test

EXPOSE 5432