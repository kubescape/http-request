FROM golang:1.18-alpine as builder

ENV RELEASE=$image_version

ENV GO111MODULE=

ENV CGO_ENABLED=0

WORKDIR /work
ADD . .
 
# build kubescape cmd
WORKDIR /work
RUN go build -o build/http_request

FROM alpine

RUN addgroup -S request && adduser -S ks -G request
USER request
WORKDIR /home/request/

COPY --from=builder /work/build/http_request /usr/bin/http_request

ENTRYPOINT ["http_request"]
