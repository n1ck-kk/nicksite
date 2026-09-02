# nicksite

Personal website built with Go, templ, and HTMX.

**Live:** https://nicksite-app.orangeriver-2bd0314b.eastus.azurecontainerapps.io

## Stack

- Go + templ + chi
- CSS (no framework)
- Docker → Azure Container Apps

## Dev

```bash
cp .env.example .env
templ generate
go run ./cmd/server
```

## Deploy

Merge to `main` — GitHub Actions builds, pushes to ACR, and deploys automatically.
