FROM golang:1.22 AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /server ./cmd/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /server /server
EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT [ "/server" ]