# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.15-buster as builder

# Create and change to the app directory.
WORKDIR /go/src/vdicalc

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY  . .
RUN go get github.com/spf13/viper
RUN go get github.com/doug-martin/goqu
RUN go get github.com/go-sql-driver/mysql
RUN go get google.golang.org/api/oauth2/v2
RUN go get github.com/google/go-cmp/cmp

# Build the binary.
RUN go build -o main .

# Expose Port
EXPOSE 8080

# Run the web service on container startup.
CMD ["./main"]
