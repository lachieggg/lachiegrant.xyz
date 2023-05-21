FROM golang:1.20-alpine

ARG BUILD_DATE
COPY .ignore .

RUN echo "$BUILD_DATE" > .ignore

WORKDIR /app

COPY go.mod go.sum ./
RUN apk update
RUN apk add procps htop aha

CMD ["./scripts/start.sh"]
