FROM golang:1.18 as builder

ENV GO111MODULE="on"
ENV CGO_ENABLED="0"
ENV GOOS="linux"

WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/pg-stat-test-task /app/cmd/

FROM gcr.io/distroless/static
COPY --from=builder /app/pg-stat-test-task /usr/bin/pg-stat-test-task
ENTRYPOINT ["/usr/bin/pg-stat-test-task"]