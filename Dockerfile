############################################
## build
############################################
FROM golang as build

WORKDIR /build
COPY . .

RUN go build -o serial-gateway cmd/serial/main.go 


############################################
## prod
############################################
FROM golang

WORKDIR /app

COPY --from=build /build/serial-gateway .
CMD ./serial-gateway

