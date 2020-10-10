# lightweight base image
FROM alpine:3.12

# adds compatibilty layer for libc between musl libc and glibc
# this is needed which go binaries compiled with 
# CGO_ENABLED=1, which is required for dependencies in this project (sqlite)
# RUN apk add --no-cache libc6-compat

# Screw it, I'm just gonna use Ubuntu for now
# TODO: Stop using Ubuntu
# FROM ubuntu:focal

# add timezone data
# this allows the user to specify a timezone as an environment variable: TZ
RUN apk add --no-cache tzdata

# copy go binary into container
COPY photolum /app/photolum

# copy db schema into container
COPY database/schema.sql /app/schema.sql

# expose port for API access
EXPOSE 8080

# set go binary as entrypoint
CMD ["/app/photolum"]