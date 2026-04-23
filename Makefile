PORT ?= 4400

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[1m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: dev
dev: ## Serve the site locally on :$(PORT)
	@echo "serving pointegrity on http://localhost:$(PORT)"
	@python3 -m http.server $(PORT)

.PHONY: check
check: ## Check for broken local links and missing assets
	@echo "checking index.html"
	@grep -oE 'href="[^"]+"|src="[^"]+"' index.html | sort -u
	@echo "---"
	@echo "checking about/index.html"
	@grep -oE 'href="[^"]+"|src="[^"]+"' about/index.html | sort -u

.PHONY: update-motif
update-motif: ## Refresh vendored motif CSS from ../motif/
	cp ../motif/tokens.css ../motif/components.css ../motif/utilities.css ../motif/responsive.css ../motif/motif.css static/motif/
	@echo "motif CSS updated — review with: git diff static/motif/"
