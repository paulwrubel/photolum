FROM alpine:latest

EXPOSE 50001

COPY photolum /app/photolum

CMD ["/app/photolum"]