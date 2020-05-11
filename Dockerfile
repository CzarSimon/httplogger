FROM golang:1.13.7-alpine3.11 AS build

# Copy source
WORKDIR /app/httplogger
COPY . .

# Download dependencies application
RUN go mod download

# Build application.
WORKDIR /app/httplogger/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:3.11 AS run

WORKDIR /opt/app
COPY --from=build /app/httplogger/cmd/cmd httplogger
ENV GIN_MODE release
CMD ["./httplogger"]