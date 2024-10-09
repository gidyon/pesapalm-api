FROM alpine
LABEL maintainer="gideon.ngk@gmail.com"
RUN apk update && \
  apk add ca-certificates && \
  update-ca-certificates && \
  rm -rf /var/cache/apk/* && \
  apk add libc6-compat && \
  apk add gcompat
EXPOSE 80 443 8080 9090
WORKDIR /app
COPY service .
ENTRYPOINT [ "/app/service" ]
CMD ["--config-file","/app/app.env","--casbin-conf","/app/casbin.conf","--static-dir","/app/out"]