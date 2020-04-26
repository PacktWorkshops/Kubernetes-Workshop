FROM golang
WORKDIR /app
COPY . .
RUN go get -d -v
RUN go build -o main .
CMD ["./main"]