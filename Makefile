# ──────────────────────────────────────────────────────────────
# URL Shortener – top-level Makefile
# ──────────────────────────────────────────────────────────────
# Service directories
AUTH_DIR        = apps/auth-service
API_DIR         = apps/url-shortener-api
REDIRECTOR_DIR  = apps/redirector-service
WEB_DIR         = apps/url-shortener-web
ANALYTICS_DIR   = apps/analytics-service
DASHBOARD_DIR   = apps/dashboard-service
COMPOSE_FILE    = ops/config/compose.yaml

# ────── Infrastructure ──────────────────────────────────────
.PHONY: infra-up infra-down

infra-up: ## Start infrastructure (Mongo, Redis, Kafka, ClickHouse)
	docker compose -f $(COMPOSE_FILE) up -d

infra-down: ## Stop infrastructure
	docker compose -f $(COMPOSE_FILE) down

# ────── Individual services ─────────────────────────────────
.PHONY: run-auth run-api run-redirector run-web run-analytics run-dashboard

run-auth: ## Run auth-service (Go, port 8081)
	cd $(AUTH_DIR) && go run .

run-api: ## Run url-shortener-api (Go, port 8080)
	cd $(API_DIR) && go run .

run-redirector: ## Run redirector-service (Go, port 8082)
	cd $(REDIRECTOR_DIR) && go run .

run-web: ## Run Angular web frontend (dev server)
	cd $(WEB_DIR) && npx ng serve

run-analytics: ## Build & run analytics-service (Flink jar)
	cd $(ANALYTICS_DIR) && ./gradlew shadowJar && java -jar build/libs/analytics-service-all.jar

run-dashboard: ## Run dashboard-service (Spring Boot)
	cd $(DASHBOARD_DIR) && ./gradlew bootRun

# ────── Convenience targets ─────────────────────────────────
.PHONY: run-all-go run-all stop-all help

run-all-go: ## Run all three Go services in separate terminals
	cmd /c start "" cmd /c "title auth-service && cd $(AUTH_DIR) && go run . & pause"
	cmd /c start "" cmd /c "title url-shortener-api && cd $(API_DIR) && go run . & pause"
	cmd /c start "" cmd /c "title redirector-service && cd $(REDIRECTOR_DIR) && go run . & pause"

run-all: infra-up ## Start infra + all services (each in its own terminal)
	cmd /c start "" cmd /c "title auth-service && cd $(AUTH_DIR) && go run . & pause"
	cmd /c start "" cmd /c "title url-shortener-api && cd $(API_DIR) && go run . & pause"
	cmd /c start "" cmd /c "title redirector-service && cd $(REDIRECTOR_DIR) && go run . & pause"
	cmd /c start "" cmd /c "title dashboard-service && cd $(DASHBOARD_DIR) && gradlew bootRun & pause"
	cmd /c start "" cmd /c "title url-shortener-web && cd $(WEB_DIR) && npx ng serve & pause"

stop-all: infra-down ## Stop infrastructure + kill service processes
	-taskkill /FI "WINDOWTITLE eq auth-service" /F 2>nul
	-taskkill /FI "WINDOWTITLE eq url-shortener-api" /F 2>nul
	-taskkill /FI "WINDOWTITLE eq redirector-service" /F 2>nul
	-taskkill /FI "WINDOWTITLE eq dashboard-service" /F 2>nul
	-taskkill /FI "WINDOWTITLE eq url-shortener-web" /F 2>nul
	@echo Infrastructure and services stopped.

# ────── Help ────────────────────────────────────────────────
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
