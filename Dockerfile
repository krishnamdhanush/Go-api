FROM golang:1.20.3

WORKDIR /app

COPY . .


RUN go build -o my-go-api

EXPOSE 8080

CMD ["./my-go-api"]