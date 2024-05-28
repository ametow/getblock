FROM golang:1.22.3 AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN make all

FROM alpine:3.18.4
ARG API_KEY
ENV GETBLOCK_KEY=$API_KEY
WORKDIR /app
COPY --from=builder /app/* /app/
COPY etc/config.yaml  /app/config.yaml

RUN apk add --no-cache libc6-compat
EXPOSE 8080

CMD [ "/app/getblock-cli", "-c", "/app/config.yaml" ]