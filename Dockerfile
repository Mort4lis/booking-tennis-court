# ---- build stage ------------------------------------------------------------
FROM golang:1.26 AS build

WORKDIR /src

# Cache module downloads.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/bot ./cmd/bot

# ---- runtime stage --------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/bot /usr/local/bin/bot

USER nonroot:nonroot

# Metrics endpoint (see BTC_METRICS_ADDR).
EXPOSE 9090

ENTRYPOINT ["/usr/local/bin/bot"]
CMD ["serve"]
