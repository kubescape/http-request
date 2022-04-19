FROM golang:1.18-alpine as builder

ENV RELEASE=$image_version

ENV GO111MODULE=

ENV CGO_ENABLED=0

WORKDIR /work
ADD . .

RUN apk add git

WORKDIR /work
RUN go build -o build/http_request

FROM alpine

RUN addgroup -S armo && adduser -S armo -G armo
USER armo
WORKDIR /home/armo/

COPY --from=builder /work/build/http_request /usr/bin/http_request

ENTRYPOINT ["http_request"]
