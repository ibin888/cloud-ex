FROM golang:alpine AS build
WORKDIR /app
COPY go.mod /app
COPY main.go /app
RUN go build .

FROM busybox
WORKDIR /bin
COPY --from=build /app/httptest /bin/httptest
ENTRYPOINT ["./httptest"]

