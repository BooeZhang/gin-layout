FROM golang:1.21.7-alpine AS gbdata
ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

COPY . .
RUN go build -o servive /app/cmd
WORKDIR /app
CMD ["/app/servive","-c","/app/etc/config.toml"]