FROM golang:1.17-alpine as buildgo

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./app

#-------

FROM alpine:latest

COPY --from=buildgo ./src/app .

EXPOSE 8080

CMD [ "./app" ]