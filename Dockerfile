FROM golang:1.23.1 AS build

WORKDIR /app

COPY . .

RUN go build -o prb_care_api cmd/prb_care_api/main.go

FROM alpine:3.20.2

RUN apk add --no-cache tzdata

ENV TZ=Asia/Jakarta

RUN apk add --no-cache libc6-compat

WORKDIR /app

COPY --from=build /app/prb_care_api /app/prb_care_api

CMD ["/app/prb_care_api"]