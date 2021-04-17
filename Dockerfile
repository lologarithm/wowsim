# Web server has a dependency on node
FROM node:latest as node
FROM golang:1.15.7

# Copy node executables to golang container
COPY --from=node /usr/local/lib/node_modules /usr/local/lib/node_modules
COPY --from=node /usr/local/bin/node /usr/local/bin/node
RUN ln -s /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm

# Set up working directory and copy everything to the container
WORKDIR /go/src/app
COPY . .

# Install dependencies
RUN go get -d -v ./...
RUN go install -v ./...

# Build server and ui
RUN GOOS=linux GOARCH=amd64 go build .
RUN GOOS=js GOARCH=wasm go build -o ui/lib.wasm ui/main_wasm.go

EXPOSE 3333

# Run the dev server
CMD go run main.go web.go --web
