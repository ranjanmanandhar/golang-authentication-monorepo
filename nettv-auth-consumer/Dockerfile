FROM golang:1.19.0-alpine3.16 as build

WORKDIR /go/src/gitlab-server.wlink.com.np/nettv/nettv-auth-consumer

COPY . .

ARG GIT_USERNAME
ARG GIT_PASSWORD

RUN apk update && apk add git

RUN go env -w GOPRIVATE=gitlab-server.wlink.com.np/*

RUN git config --global url.'https://'$GIT_USERNAME:$GIT_PASSWORD'@gitlab-server.wlink.com.np'.insteadOf \
    "https://gitlab-server.wlink.com.np"    

RUN go build -o app

FROM alpine:3.16

COPY --from=build /go/src/gitlab-server.wlink.com.np/nettv/nettv-auth-consumer/app /usr/local/bin/nettv-auth-consumer
COPY ./bin/wait-for-it.sh /
RUN chmod +x /wait-for-it.sh
CMD ["/usr/local/bin/nettv-auth-consumer"]
