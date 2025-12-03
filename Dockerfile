FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go mod init proxy
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o proxy .

FROM gcr.io/distroless/static:nonroot
COPY --from=build /app/proxy /proxy
EXPOSE 8080
ENTRYPOINT ["/proxy"]
