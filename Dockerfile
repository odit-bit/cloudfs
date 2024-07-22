FROM golang:1.22.5-alpine3.20 as build-stage

# WORKDIR /usr/src/app

# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

# COPY . .
# RUN go build  -v -o ./web ./cmd/web

# # distroless
# FROM gcr.io/distroless/base-debian12

# WORKDIR /app

# COPY --from=build-stage /usr/src/app/web .

# CMD [ "./web" ]


WORKDIR /usr/src/app

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x


# Build the application.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /cloudfs-web ./cmd/app

# default config file
COPY cloudfs.yaml .

# distroless
FROM gcr.io/distroless/base-debian12 as Final-stage
WORKDIR /app

COPY cloudfs.yaml .
COPY --from=build-stage /cloudfs-web .

EXPOSE 8181
ENTRYPOINT ["./cloudfs-web"]
CMD [ "run" ]