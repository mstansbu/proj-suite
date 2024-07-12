#build stage
FROM golang:1.22 AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v ./...

#final stage
FROM alpine:latest
COPY --from=builder /usr/src/app /app
LABEL Name=projsuite Version=0.0.1
EXPOSE 3000
CMD [ "./tic-tac-toe" ]