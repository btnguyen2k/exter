## Sample Dockerfile to package Frontend-GUI as a docker image.
# Sample build command:
# docker build --force-rm --squash -t exter_fe .

FROM node:12.22-alpine3.15 AS builder
RUN mkdir /build
COPY . /build
RUN cd /build && rm -rf dist node_modules && npm install && npm run build

FROM nginx:1.17-alpine
LABEL maintainer="Thanh Nguyen <btnguyen2k@gmail.com>"
COPY nginx.conf /etc/nginx/nginx.conf
RUN mkdir -p /usr/share/nginx/html/app
COPY --from=builder /build/dist /usr/share/nginx/html/app
COPY --from=builder /build/dist/favicon*.* /usr/share/nginx/html/
COPY --from=builder /build/dist/manifest.json /usr/share/nginx/html
COPY index.html /usr/share/nginx/html
EXPOSE 80
RUN chown nginx.nginx /usr/share/nginx/html/ -R
