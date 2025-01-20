FROM docker.io/golang as build

WORKDIR /app
COPY . /app
RUN go mod tidy
WORKDIR /app/cmd
RUN go build -o main

FROM docker.io/golang

COPY --from=build ./app/cmd/main .
CMD ["./main"]
