FROM golang:alpine as build-stage

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build  -v -o ./web ./cmd/web

# distroless
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=build-stage /usr/src/app/web .

CMD [ "./web" ]