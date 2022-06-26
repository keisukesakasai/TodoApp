FROM golang:1.17.6 as builder
WORKDIR /workspace
COPY . /workspace
RUN CGO_ENABLED=0 GOOS=linux go build -o todoapp main.go && chmod +x ./todoapp

FROM alpine:3.15
WORKDIR /app
RUN apk --no-cache add ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /workspace/ ./

EXPOSE 8888
CMD ["./todoapp"]