# Encrypted Vault

## env vars

VAULT_MODE=kms|local
VAULT_BUCKET=<s3-bucket> (required for kms)
VAULT_KMS_KEY=<kms-key-id> (required for kms)
VAULT_REMOTE_PATH=/path/to/local/vault (required for local mode)
VAULT_REQUIRE_PASSWORD=1 (optional)
VAULT_PASS_FILE=/path/vault_pass.txt
VAULT_DB_PATH=vault.db

## build & run

go mod tidy
go build -o vault
./vault upload secret.txt

## docker

docker build -t vault-cli .
docker run --rm -e AWS_ACCESS_KEY_ID=... -e AWS_SECRET_ACCESS_KEY=... -e VAULT_BUCKET=... -e VAULT_KMS_KEY=... vault-cli upload file.txt
