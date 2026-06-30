EVENTS_PROTO_DIR ?= proto/events
EVENTS_OUT ?= gen/events/tenant/v1

# =============================================================================
# Events (buf + protoc-gen-go)
# =============================================================================

.PHONY: events-generate
events-generate: _connect-check-tools ## Generate Go types from Kafka event proto files
	@echo "$(COLOR_BLUE)→ Generating event types...$(COLOR_RESET)"
	$(BUF) generate $(EVENTS_PROTO_DIR) --template buf.gen.events.yaml
	@echo "$(COLOR_GREEN)✓ Event generation complete!$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)  Generated files in $(EVENTS_OUT)/:$(COLOR_RESET)"
	@find $(EVENTS_OUT) -name '*.pb.go' | sed 's/^/    /'

.PHONY: events-clean
events-clean: ## Remove generated event pb.go files
	@echo "$(COLOR_BLUE)→ Cleaning generated event files...$(COLOR_RESET)"
	@find $(EVENTS_OUT) -name '*.pb.go' -delete
	@echo "$(COLOR_GREEN)✓ Cleaned event pb.go files$(COLOR_RESET)"
