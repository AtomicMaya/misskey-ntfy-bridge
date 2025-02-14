# Build stage 
FROM golang:1.23.4 AS builder

RUN apt update && apt upgrade -y

WORKDIR /app

COPY . .

RUN mkdir -p ./build

RUN GCO_ENABLED=0 go build -tags netgo -ldflags "-w" -a -o ./build/misskey-ntfy-bridge-latest ./app

# Final stage
FROM alpine

RUN mkdir -p /release

COPY --from=builder /app/build/ /release

ENV SOURCE=container
ENV HOST=0.0.0.0
ENV PORT=1337

WORKDIR /release

RUN chmod +x ./misskey-ntfy-bridge-latest

CMD ["/release/misskey-ntfy-bridge-latest"]