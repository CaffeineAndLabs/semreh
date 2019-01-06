FROM golang:1.11.3-alpine AS builder

RUN apk add --update --no-cache git ca-certificates tzdata zip

WORKDIR /src
COPY . ./

RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
    CGO_ENABLED=0 go build \
    -mod=vendor \
    -installsuffix 'static' \
    -ldflags "-X main.GitCommit=$GIT_COMMIT" \
    -o /app .

WORKDIR /usr/share/zoneinfo
# go tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM scratch AS final

# the timezone data:
ENV ZONEINFO /zoneinfo.zip
COPY --from=builder /zoneinfo.zip /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

EXPOSE 8000

ENTRYPOINT ["/app"]
