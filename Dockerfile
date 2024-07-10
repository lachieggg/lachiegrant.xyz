FROM golang:1.22-alpine

ARG BUILD_DATE
COPY .ignore .

RUN echo "$BUILD_DATE" > .ignore

WORKDIR /app

COPY go.mod go.sum ./
RUN apk update
RUN apk add procps htop aha neofetch

CMD ["./scripts/start.sh"]
