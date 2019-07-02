FROM golang:1.12-alpine3.10 AS build
RUN apk update && apk add git

# Copy source
WORKDIR /app/httplogger
COPY . .

# Build application
WORKDIR /app/httplogger/cmd
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:3.10 AS run

WORKDIR /opt/app
COPY --from=build /app/httplogger/cmd/cmd httplogger
ENV GIN_MODE release
CMD ["./httplogger"]