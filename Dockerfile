# Start from the golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Install Delve for remote debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

RUN go mod tidy

# Download all the dependencies
RUN go mod download


RUN go get github.com/microsoft/go-mssqldb/integratedauth/ntlm@v1.6.0

# Expose port 40000 for the debugger
EXPOSE 40000



# Command to run the application with Delve
#CMD ["dlv", "debug", ".", "--headless", "--listen=:40000", "--log", "--api-version=2"]

CMD ["tail", "-f", "/dev/null"]

#dlv debug . --headless --listen=:40000 --log --api-version=2