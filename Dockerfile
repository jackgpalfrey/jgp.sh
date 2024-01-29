FROM golang:1.21-alpine

WORKDIR /app

COPY src/go.mod .
RUN go mod download
COPY src/ .
RUN go build -o /app-executable

EXPOSE 8080
CMD [ "/app-executable" ]
