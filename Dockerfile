## Multistage build
FROM golang:1.23.1 as build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o /app

## Multistage deploy
FROM gcr.io/distroless/base-debian12

WORKDIR /
COPY --from=build /src/config /config
COPY --from=build /src/template /template
COPY --from=build /app /app

ENTRYPOINT ["/app"]
