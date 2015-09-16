FROM alpine:3.2

RUN echo 'http://dl-4.alpinelinux.org/alpine/edge/main' >> /etc/apk/repositories \
    && apk update && apk add go git ca-certificates make && rm -rf /var/cache/apk/* \
    && wget -P /usr/local/bin http://public.thisissoon.com.s3.amazonaws.com/glide \
    && chmod +x /usr/local/bin/glide

ENV GOPATH=/deadpool

RUN go get github.com/gorilla/mux
EXPOSE 80

WORKDIR /deadpool/src/github.com/thisissoon/deadpool
COPY . /deadpool/src/github.com/thisissoon/deadpool

RUN make install
RUN ln -s /deadpool/bin/deadpool /usr/local/bin/deadpool

ENTRYPOINT deadpool
