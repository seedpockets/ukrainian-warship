FROM --platform=linux/x86-64 golang:1.17
# FROM golang:1.17
ENV GO111MODULE=on
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app
# RUN go build -mod=vendor -o invoice
RUN go build -ldflags "-s -w" -o ukrainian-warship

# Command to run the executable
CMD ["/app/ukrainian-warship", "kill"]