FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
ARG SERVICE
RUN CGO_ENABLED=0 go build -o /app ./${SERVICE}

FROM alpine:3.20
COPY --from=build /app /app
ENTRYPOINT ["/app"]
