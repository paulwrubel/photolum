# lightweight base image
# FROM alpine:3.12

# adds compatibilty layer for libc between musl libc and glibc
# this is needed which go binaries compiled with 
# CGO_ENABLED=1, which is required for dependencies in this project (sqlite)
# RUN apk add --no-cache libc6-compat

# Screw it, I'm just gonna use Ubuntu for now
# TODO: Stop using Ubuntu
FROM ubuntu:bionic

# copy go binary into container
COPY photolum /app/photolum

# expose port for API access
EXPOSE 8080

# set go binary as entrypoint
CMD ["/app/photolum"]