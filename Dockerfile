#FROM golang:alpine as builder
FROM --platform=linux/x86-64 golang:alpine as builder
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w" -o ukrainian-warship

FROM --platform=linux/x86-64 alpine:3.14
WORKDIR /app

COPY --from=builder /app/ukrainian-warship ./
COPY --from=builder /app/default_targets.json ./

# Command to run the executable
CMD ["/app/ukrainian-warship", "kill", "--workers=24"]