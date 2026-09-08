# ---- build ----
FROM golang:1.22-alpine AS build
WORKDIR /src

# go.sum se genera en el build para no versionar hashes en el reto.
COPY go.mod ./
COPY . .
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/qr-api .

# ---- test (opcional: docker build --target test) ----
FROM build AS test
RUN go test ./... -v

# ---- runtime ----
FROM alpine:3.20
RUN adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=build /out/qr-api .
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/qr-api"]
