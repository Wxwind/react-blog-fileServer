FROM golang:1.19.5
ENV GOPROXY https://goproxy.cn
RUN mkdir /app
WORKDIR /app
COPY . /app
RUN go mod tidy
RUN go build -o ./fileServerEXE ./main.go
ENV GIN_MODE=release
CMD ["./fileServerEXE"]