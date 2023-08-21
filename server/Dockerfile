FROM golang:alpine as preprod
WORKDIR /app

COPY  go.mod go.sum ./
RUN go mod download

COPY . .
RUN apk update && apk add ffmpeg

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main

CMD ["/main"]