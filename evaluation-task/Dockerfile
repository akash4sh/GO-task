FROM golang:alpine
RUN mkdir /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /evalution

EXPOSE 8000

CMD [ "/evalution" ]