# URL Shortener Monorepo

This repository is organized as a small monorepo with each runnable application under `apps/` and shared operational configuration under `ops/`.

## Structure

- `apps/url-shortener-api`: Go API for creating and resolving short URLs
- `apps/auth-service`: Go auth service that issues JWTs
- `apps/url-shortener-web`: Angular frontend
- `apps/analytics-service`: analytics service
- `apps/dashboard-service`: dashboard service
- `ops/config`: shared deployment or compose configuration

## Notes

- Build each app from its own directory.
- Repo-level ignore rules are scoped to the new monorepo layout.
