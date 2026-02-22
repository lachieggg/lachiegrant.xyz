FROM golang:1.26-alpine

# BUILD_DATE used here to invalidate cache for
# subsequent layers when it changes.
ARG BUILD_DATE
RUN echo "Build Date: $BUILD_DATE" > /dev/null

WORKDIR /app

COPY . .
RUN go build -o /app/bin/app ./src

CMD ["/app/bin/app"]
