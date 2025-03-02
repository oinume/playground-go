atlas migrate apply \
  --dir "file://migrations" \
  --url "postgres://postgres:pass@localhost:5432/atlas-db?search_path=public&sslmode=disable"
