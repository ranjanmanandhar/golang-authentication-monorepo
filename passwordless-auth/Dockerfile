FROM golang:1.19.0 as build

WORKDIR /go/src/gitlab-server.wlink.com.np/nettv/nettv-auth/passwordless-auth

COPY . .

ARG GIT_USERNAME
ARG GIT_PASSWORD

# RUN apk update && apk add git

# RUN apk add build-base

RUN go env -w GOPRIVATE=gitlab-server.wlink.com.np/*

RUN git config --global url.'https://'$GIT_USERNAME:$GIT_PASSWORD'@gitlab-server.wlink.com.np'.insteadOf \
    "https://gitlab-server.wlink.com.np"

ENV CGO_enabled=1

RUN go mod tidy

RUN go build -o app

FROM gitlab-server.wlink.com.np:4567/nettv/nettv-auth/go/ubuntu:oci8

COPY --from=build /go/src/gitlab-server.wlink.com.np/nettv/nettv-auth/passwordless-auth/app /usr/local/bin/passwordless

CMD ["/usr/local/bin/passwordless"]