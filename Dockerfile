FROM golang:1.14

# Files setup
WORKDIR /go/src/k-cad
COPY . .

ENV PORT 3000
RUN go get -d -v ./...
RUN go install -v ./...

CMD go run kcad.go -port=${PORT}