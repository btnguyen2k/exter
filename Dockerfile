## Sample Dockerfile to package the whole application (backend and frontend) as a Docker image.
# Sample build command:
# docker build --rm -t exter .

FROM node:13.6-alpine3.11 AS builder_fe
RUN apk add jq
RUN mkdir /build
COPY ./fe-gui /build
COPY ./fe-gui/src/config_one_image.json /build/src/config.json
RUN cd /build && rm -rf dist node_modules && export BUILD=`date +%Y%m%d%H%M` && sed -i 's/exter - Exter v0.1.0/GoVueAdmin v0.1.1 b'$BUILD/g public/index.html && npm install && npm run build

FROM golang:1.13-alpine AS builder_be
RUN apk add git build-base jq
RUN mkdir /build
COPY ./be-api /build
RUN cd /build && go build -o main

FROM alpine:3.10
LABEL maintainer="Thanh Nguyen <btnguyen2k@gmail.com>"
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
