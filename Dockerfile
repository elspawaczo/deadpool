# Ubuntu just works
FROM alpine:3.2
MAINTAINER SOON_ <dorks@thisissoon.com>

## Environment Variables
ENV GOPATH /deadpool
ENV GOBIN /usr/local/bin
ENV PATH $PATH:$GOPATH/bin

# OS Dependencies
RUN echo 'http://dl-4.alpinelinux.org/alpine/edge/main' >> /etc/apk/repositories \
    && apk update && apk add go go-tools git gcc g++ ca-certificates make bash && rm -rf /var/cache/apk/*

# Set working Directory
WORKDIR /deadpool

# GPM (Go Package Manager)
RUN git clone https://github.com/pote/gpm.git \
    && cd gpm \
    && git checkout v1.3.2 \
    && ./configure \
    && make install

# Install Dependencies
COPY ./Godeps /deadpool/Godeps
RUN gpm install && go get upper.io/db/postgresql

# Set our final working dir to be where the source code lives
WORKDIR /deadpool/src/github.com/thisissoon/deadpool

# Set the default entrypoint to be deadpool
ENTRYPOINT ["deadpool"]

# Copy source code into the deadpool src directory so Go can build the package
COPY . /deadpool/src/github.com/thisissoon/deadpool

# Install the go package
RUN go install
