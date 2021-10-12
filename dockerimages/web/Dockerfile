# Telling to use Docker's golang ready image
FROM golang

MAINTAINER Guillermo de la Puente <g.puente.allott@gmail.com>

EXPOSE 8080 80 443

# Set pwd to the go folder
WORKDIR /app

RUN apt-get update -y && \
	apt-get install -y \
	vim
RUN go get "github.com/go-sql-driver/mysql"
# Install and build binaries revel for this server
# Copy the entrypoint file

# Start process
#ENTRYPOINT ["go mod download","go run main.go"]
