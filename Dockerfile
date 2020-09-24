# lightweight base image
FROM alpine:3.12

# copy go binary into container
COPY photolum /app/photolum

# expose port for API access
EXPOSE 50001

# set go binary as entrypoint
CMD ["/app/photolum"]