FROM golang:1.11.3-alpine AS builder

RUN apk add --update --no-cache git ca-certificates

WORKDIR /src
COPY . ./

RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
    CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -ldflags "-X main.GitCommit=$GIT_COMMIT" \
    -o /app .

FROM scratch AS final
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

EXPOSE 8000

ENTRYPOINT ["/app"]
