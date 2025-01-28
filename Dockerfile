FROM ubuntu:latest
ENV DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -y ca-certificates && update-ca-certificates
COPY platform-take-home /usr/bin/platform-take-home
RUN chmod a+x /usr/bin/platform-take-home
EXPOSE 8080 8081 9008
ENTRYPOINT ["/usr/bin/platform-take-home"]

