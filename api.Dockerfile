# build stage
FROM golang:1.24-alpine AS build

# copy & build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -v -o /usr/local/bin/app ./api/cmd

# prodcution stage
FROM alpine:latest

COPY --from=build /usr/local/bin/app /usr/local/bin/app

EXPOSE 8080

CMD ["app"]

