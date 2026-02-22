FROM golang:1.26-alpine

ARG BUILD_DATE
# BUILD_DATE is used here to invalidate cache for subsequent layers when it changes.
RUN echo "Build Date: $BUILD_DATE" > /dev/null

WORKDIR /app

COPY . .
RUN go build -o /app/bin/app ./src

CMD ["/app/bin/app"]
