FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app ./cmd/todo/main.go

FROM alpine:latest AS app
WORKDIR /app
COPY --from=builder /app/todo-app ./
COPY --from=builder /app/web ./web
RUN mkdir /data
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
EXPOSE ${TODO_PORT} 
ENTRYPOINT ["./todo-app"]