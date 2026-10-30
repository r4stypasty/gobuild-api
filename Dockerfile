FROM golang:1.27.3-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /out/gobuild-api \
	./cmd/api

FROM dh-mirror.gitverse.ru/library/alpine:3.23

RUN addgroup -S app && adduser -S -G app app && mkdir -p /data && chown app:app /data

COPY --from=build /out/gobuild-api /usr/local/bin/gobuild-api

WORKDIR /data
VOLUME ["/data"]
EXPOSE 8080
USER app

ENTRYPOINT ["/usr/local/bin/gobuild-api"]
