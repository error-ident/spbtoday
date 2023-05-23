FROM golang:1.19 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 go build -o /app/spbtoday main.go


# Deploy
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/spbtoday .
EXPOSE 1337
CMD ["/app/spbtoday"]