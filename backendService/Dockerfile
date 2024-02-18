# Use an official Golang runtime as a parent image
FROM golang:1.21.3

# Set the working directory inside the container
WORKDIR /app
# Create a directory named "MessageSystem" in the container

# Copy the current directory contents into the container at /app
COPY . .




# Build the Go application inside the container
RUN go build -o main .



# Define the command to run your Go server when the container Starts
CMD ["./main"]
