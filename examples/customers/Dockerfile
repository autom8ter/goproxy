FROM golang:alpine AS builder
COPY ./main.go /go/src/github.com/autom8ter/customers/main.go
COPY ./vendor /go/src/github.com/autom8ter/customers/vendor

RUN set -ex && \
  cd /go/src/github.com/autom8ter/customers && \
  CGO_ENABLED=0 go build \
        -tags netgo \
        -v -a \
        -ldflags '-extldflags "-static"' && \
  mv ./customers /usr/bin/customers

FROM busybox
ENV SECRET=somesecret
# Retrieve the binary from the previous stage
COPY --from=builder /usr/bin/customers /usr/local/bin/customers

# Set the binary as the entrypoint of the container
ENTRYPOINT [ "customers" ]