FROM golang:1.19-bullseye AS build

# Copy source
WORKDIR /app/httplogger
COPY go.mod .
COPY go.sum .

# Download dependencies application
RUN go mod download

# Build application.
COPY cmd cmd
COPY internal internal
WORKDIR /app/httplogger/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o httplogger

FROM gcr.io/distroless/base-debian11 AS runtime

# Copy binary from buid step
WORKDIR /opt/app
COPY --from=build /app/httplogger/cmd/httplogger httplogger

# Prepare runtime
USER nonroot:nonroot
ENV GIN_MODE release
CMD ["./httplogger"]