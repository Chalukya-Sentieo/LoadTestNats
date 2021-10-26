FROM golang:1.17.2 as builder

WORKDIR /app/
COPY ../src/go.mod ../src/go.sum /app/
RUN go mod tidy
COPY ../src/ .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /GinServer

# final stage
FROM alpine
ARG PARALLEL_BATCH=10
ARG SEQUENTIAL_BATCH=100
ARG NATS_SERVERS="nats://10.0.63.21:4222,nats://10.0.63.20:4222,nats://10.0.63.161:4222,nats://10.0.63.174:4222,nats://10.0.63.214:4222"
#ARG NATS_SERVERS="nats://13.233.167.108:4222,nats://13.127.160.94:4222,nats://13.234.110.245:4222,nats://13.235.245.159:4222,nats://52.66.214.53:4222"
ENV NATS_SERVERS="${NATS_SERVERS}"
ENV PARALLEL_BATCH="${PARALLEL_BATCH}"
ENV SEQUENTIAL_BATCH="${SEQUENTIAL_BATCH}"
COPY --from=builder /GinServer /GinServer
EXPOSE 9999
ENTRYPOINT ["/GinServer"]