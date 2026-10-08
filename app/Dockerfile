# syntax=docker/dockerfile:1
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# distroless/static: không shell, chạy non-root
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server
USER 10001:10001
EXPOSE 8080
ENTRYPOINT ["/server"]
