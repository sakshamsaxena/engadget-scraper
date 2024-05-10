FROM golang:1.19

WORKDIR /app
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o my-app
ARG mode=development
ENV APP_ENV=$mode
CMD ["./my-app"]