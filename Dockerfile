# build stage
FROM golang:1.25-alpine AS build

WORKDIR /build

RUN mkdir -p /etc/corsair/ && touch /etc/corsair/config.yaml

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /corsair ./cmd


# final image
FROM gcr.io/distroless/static AS build-release-stage

USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /corsair /corsair
COPY --from=build --chown=nonroot:nonroot /etc/corsair /etc/corsair

ENTRYPOINT ["/corsair", "-c", "/etc/corsair/config.yaml"]