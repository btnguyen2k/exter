## Sample Dockerfile to package the whole application (backend and frontend) a docker image.
# Sample build command:
# docker build --force-rm --squash -t exter:0.1.0 .

FROM node:13.6-alpine3.11 AS builder_fe
RUN mkdir /build
COPY ./fe-gui /build
COPY ./fe-gui/src/config_one_image.json /build/src/config.json
RUN cd /build && npm install && npm run build

FROM golang:1.13-alpine AS builder_be
RUN apk add git build-base
RUN mkdir /build
COPY ./be-api /build
RUN cd /build && go build -o main

FROM alpine:3.10
MAINTAINER Thanh Nguyen <btnguyen2k@gmail.com>
RUN mkdir /app
RUN mkdir /app/frontend
COPY --from=builder_be /build/main /app/main
COPY --from=builder_be /build/config /app/config
COPY --from=builder_fe /build/dist /app/frontend
RUN apk add --no-cache -U tzdata bash ca-certificates \
    && update-ca-certificates \
    && cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime \
    && chmod 711 /app/main \
    && rm -rf /var/cache/apk/*
WORKDIR /app
EXPOSE 8000
CMD ["/app/main"]
#ENTRYPOINT /app/main
