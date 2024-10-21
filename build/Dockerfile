# Build the application from source
FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY ./pkg/ ./pkg/
COPY ./cmd ./cmd

RUN find ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /tf-profile

# Run the tests in the container
FROM build AS run-test-stage
COPY ./test ./test
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build /tf-profile /tf-profile

USER nonroot:nonroot

ENTRYPOINT ["/tf-profile"]