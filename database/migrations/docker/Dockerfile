FROM alpine:3.13

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ENV MIGRATION_ROOT=database/migrations
ENV MIGRATION_DIR=$MIGRATION_ROOT/docker
ADD https://github.com/pressly/goose/releases/download/v3.22.1/goose_linux_x86_64 bin/goose
RUN chmod +x bin/goose

WORKDIR /root

ADD $MIGRATION_ROOT/*.sql .
ADD $MIGRATION_DIR/migration.sh .
ADD .env .
ADD bin/goose .

RUN chmod +x migration.sh

#ENTRYPOINT ["bash", "./migration.sh"]
