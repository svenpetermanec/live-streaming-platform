# build stage
FROM golang:1.24-alpine AS build

# compilation dependecies
RUN apk add --no-cache \
    build-base \
    tcl \
    cmake \
    pkgconfig \
    openssl-dev \
    curl \
    linux-headers

# SRT dependencies
WORKDIR /srt
RUN wget -O srt.tar.gz https://github.com/Haivision/srt/archive/refs/tags/v1.5.3.tar.gz && \
    tar xf srt.tar.gz && \
    cd srt-1.5.3 && \
    ./configure && \
    make && \
    make install

# copy & build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -v -o /usr/local/bin/app ./transcoder/cmd

# production stage
FROM alpine:latest

RUN apk add --no-cache \
    ffmpeg \
    openssl \
    ca-certificates

COPY --from=build /usr/local/lib/libsrt.so* /usr/local/lib
COPY --from=build /usr/local/include/srt /usr/local/include/srt

COPY --from=build /usr/local/bin/app /usr/local/bin/app

RUN ldconfig /usr/local/lib

EXPOSE 5270

CMD ["app"]