FROM golang:1.23-alpine AS builder

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

WORKDIR /app/cmd/
RUN go get -v ./... \
&& go install -v ./... \
&& go build -v -o mail 

# EXPOSE 8080

FROM scratch AS production

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cmd/mail /app/mail

# HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
# CMD ["CMD-SHELL", "curl --fail http://localhost:4400/api/v1/health | grep ok && echo 0 || echo 1"]

ENTRYPOINT [ "/app/mail" ]
