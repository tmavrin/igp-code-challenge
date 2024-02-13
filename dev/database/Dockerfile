FROM postgres:16

COPY *.sql /docker-entrypoint-initdb.d/

ENV POSTGRES_DB=main
ENV POSTGRES_USER=dev
ENV POSTGRES_PASSWORD=test

EXPOSE 5432