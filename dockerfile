FROM golang:1.22-alpine

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o vault

ENV VAULT_BUCKET=mybucket \
    VAULT_KMS_KEY=mykeyid \
    VAULT_REMOTE_PATH=/vault_storage

CMD ["./vault"]