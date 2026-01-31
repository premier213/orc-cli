.PHONY: help rathole-install rathole-run scanner-port scanner-build

# Default ports to scan
DEFAULT_PORTS := 80 443 1080 2080 2081 2096 8000 8080 8081 8880 8881

# Default target - show help
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "ORC CLI - Available commands:"
	@echo ""
	@echo "  make rathole-install    Install rathole (requires root)"
	@echo "  make rathole-run        Run rathole (requires root)"
	@echo "  make scanner-build      Build the port scanner binary"
	@echo "  make scanner-port <IP> [ports]  Run port scanner"
	@echo "                         If no ports specified, uses default ports:"
	@echo "                         $(DEFAULT_PORTS)"
	@echo "                         Example: make scanner-port 1.1.1.1 8080 8081"
	@echo "  make help              Show this help message"
	@echo ""

rathole-install: ## Install rathole
	@if [ "$(shell id -u)" -ne 0 ]; then \
		echo "Error: rathole-install requires root privileges. Use: sudo make rathole-install"; \
		exit 1; \
	fi
	@bash rathole/install.sh

rathole-run: ## Run rathole
	@if [ "$(shell id -u)" -ne 0 ]; then \
		echo "Error: rathole-run requires root privileges. Use: sudo make rathole-run"; \
		exit 1; \
	fi
	@bash rathole/run.sh

scanner-build: ## Build the port scanner binary
	@cd scanner-port && go build -o scanner-port main.go
	@echo "Scanner built successfully: scanner-port/scanner-port"

scanner-port: ## Run port scanner (make scanner-port <IP> [ports])
	@ARGS="$(filter-out $@,$(MAKECMDGOALS))"; \
	if [ -z "$$ARGS" ]; then \
		echo "Error: IP address is required"; \
		echo "Usage: make scanner-port <IP> [port1] [port2] ..."; \
		echo "Example: make scanner-port 1.1.1.1 8080 8081"; \
		exit 1; \
	fi; \
	ARG_COUNT=$$(echo $$ARGS | wc -w); \
	if [ $$ARG_COUNT -eq 1 ]; then \
		echo "No ports specified, using default ports: $(DEFAULT_PORTS)"; \
		./scanner-port/scanner-port $$ARGS $(DEFAULT_PORTS); \
	else \
		./scanner-port/scanner-port $$ARGS; \
	fi

# Catch-all target to prevent Make from complaining about unknown targets (must be last)
%:
	@:
