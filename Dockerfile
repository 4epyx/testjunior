FROM golang:1.21.0-alpine AS builder

WORKDIR /build
COPY . .
RUN go build -o app cmd/app.go


FROM alpine

WORKDIR /server
COPY  --from=builder /build/app app
CMD [ "./app" ]
