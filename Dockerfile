############################################
## build
############################################
FROM golang as build

WORKDIR /build
COPY . .

RUN go build -o malone-sidecar cmd/serial/main.go 


############################################
## prod
############################################
FROM golang

WORKDIR /app

COPY --from=build /build/malone-sidecar .
CMD ./malone-sidecar

