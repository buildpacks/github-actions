FROM golang:1.20 as build-stage

WORKDIR /src
ENV GO111MODULE=on CGO_ENABLED=0

COPY . .

ARG SOURCE
RUN go build \
  -trimpath \
  -ldflags "-s -w -extldflags '-static'" \
  -tags netgo \
  -o /bin/action \
  $SOURCE/main.go

RUN strip /bin/action

FROM scratch
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-stage /bin/action /bin/action
ENTRYPOINT ["/bin/action"]
