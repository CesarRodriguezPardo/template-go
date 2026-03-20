FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  -buildvcs=false -o api-grd .

FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/api-grd .
COPY --from=builder /app/mailer/templates ./mailer/templates
COPY --from=builder /app/assets ./assets

CMD ["./api-grd"]