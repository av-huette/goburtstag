FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go build main.go

RUN ls -la .

RUN chmod +x main

CMD ["./main"]
