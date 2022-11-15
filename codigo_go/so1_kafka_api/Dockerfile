FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# RUN go build -o /docker-gs-ping

RUN go build -o api   *.go


EXPOSE 80
CMD ["ls"]
CMD [ "./api" ]