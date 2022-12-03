FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o ./bin/sybo ./cmd


FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /app/bin/sybo sybo

EXPOSE 8080

ENTRYPOINT [ "./sybo" ]
