# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./iam-service/go.mod ./iam-service/go.sum ./

# Copy the local dependency
COPY ./oort ../oort
COPY ./magnetar ../magnetar

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY ./oort ../oort
COPY ./magnetar ../magnetar

# Copy everything from the current directory to the Working Directory inside the container
COPY ./iam-service/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8001 to the outside world
EXPOSE 8001

# Command to run the executable
CMD ["./main"] 

#FROM cassandra:latest
#COPY ./iam-service/model/db/create_schema.cql ./create_schema.cql
#ENV CASSANDRA_BROADCAST_ADDRESS="cassandra"
#CMD ["cqlsh", "cassandra", "-f", "/create_schema.cql"]


