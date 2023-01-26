FROM golang:latest
WORKDIR /app
COPY . .
RUN go get -v 
RUN go build -o main .
EXPOSE 9092
EXPOSE 2181
CMD ["./main"]