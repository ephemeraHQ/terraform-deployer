ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}-alpine
WORKDIR /app
RUN apk add --no-cache git
RUN git config --system --add safe.directory '*'
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /usr/bin/terraform-deployer cmd/deploy/main.go
ENTRYPOINT /usr/bin/terraform-deployer
