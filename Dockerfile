# syntax=docker/dockerfile:1
FROM golang:alpine3.20 as build-stage


WORKDIR /usr/src/app

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x


# Build the application.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /cfs ./cmd/web

# default config file

# distroless
FROM gcr.io/distroless/base-debian12 as Final-stage
WORKDIR /app

COPY --from=build-stage /cfs .

EXPOSE 8181
ENTRYPOINT ["./cfs"]