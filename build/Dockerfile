FROM golang:alpine as builder
LABEL maintainer="Will Fantom <willfantom@gmail.com>"
RUN apk add --no-cache build-base
COPY ./go.mod /app/go.mod
WORKDIR /app
RUN go mod download
COPY . /app/
WORKDIR /app
RUN CGO_ENABLED=1 go build -o lu-covid-api .

FROM alpine
COPY --from=builder /app/database /app/database
COPY --from=builder /app/lu-covid-api /app/
EXPOSE 8080
WORKDIR /app
ENTRYPOINT [ "./lu-covid-api" ]
