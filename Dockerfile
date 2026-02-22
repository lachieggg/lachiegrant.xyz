FROM golang:1.26-alpine
ARG ALPINE_REPO

RUN apk add --no-cache htop aha --repository=$ALPINE_REPO

WORKDIR /app

COPY . .
RUN go build -o /app/bin/app ./src

CMD ["/app/bin/app"]
