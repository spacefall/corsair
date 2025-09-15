# build stage
FROM golang:1.25-alpine AS build

WORKDIR /build

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev

#RUN mkdir -p /etc/corsair/ && touch /etc/corsair/config.yaml
RUN touch config.yaml

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o corsair -ldflags "-s -w -X main.version=${VERSION}" ./cmd/

# final image
FROM gcr.io/distroless/static AS build-release-stage

USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /build/corsair /corsair
COPY --from=build --chown=nonroot:nonroot /build/config.yaml /etc/corsair/

ENTRYPOINT ["/corsair"]