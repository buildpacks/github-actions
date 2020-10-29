FROM golang:1.15 as build-stage

RUN apt-get update && apt-get install -y --no-install-recommends upx

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

RUN upx -q -9 /bin/action

FROM scratch
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-stage /bin/action /bin/action
ENTRYPOINT ["/bin/action"]
