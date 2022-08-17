# FROM golang:1.18-alpine as development
# # Add a work directory
# WORKDIR /alarmservice
# # Cache and install dependencies
# COPY go.mod go.sum ./
# RUN go mod download
# # Copy alarmservice files
# COPY . .
# # Install Reflex for development
# RUN go install github.com/cespare/reflex@latest
# # Expose port
# EXPOSE 9000
# # Start alarmservice
# CMD reflex -g '*.go' go run main.go --start-service

# FROM golang:1.18 as builder
# # Define build env
# ENV GOOS linux
# ENV CGO_ENABLED 0
# # Add a work directory
# WORKDIR /alarmservice
# # Cache and install dependencies
# COPY go.mod go.sum ./
# RUN go mod download
# # Copy alarmservice files
# COPY . .
# # Build alarmservice
# RUN go build -o alarmservice

# FROM alpine:3.14 as production
# # Add certificates
# RUN apk add --no-cache ca-certificates
# # Copy built binary from builder
# COPY --from=builder alarmservice .
# # Expose port
# EXPOSE 9000
# # Exec built binary
# CMD ./alarmservice

FROM golang:1.18-alpine AS development

ENV PROJECT_PATH=/alarmservice
ENV PATH=$PATH:$PROJECT_PATH/build
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

RUN apk add --no-cache ca-certificates tzdata make git bash

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN go mod download
RUN make

FROM alpine:3.11.2 AS production

RUN apk --no-cache add ca-certificates tzdata
COPY --from=development /alarmservice/build/alarmservice /usr/bin/alarmservice
ENTRYPOINT ["/usr/bin/alarmservice"]