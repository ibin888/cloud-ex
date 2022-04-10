FROM golang:alpine AS build
WORKDIR /app
COPY go.mod /app
COPY main.go /app
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN CGO_ENABLED=0 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM busybox
WORKDIR /bin
COPY --from=build /app/httptest /bin/httptest
ENTRYPOINT ["./httptest"]
