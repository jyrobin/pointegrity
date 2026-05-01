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

# --- Letterbox (newsletter subscriber capture) ----------------------------
# See deploy/SETUP-letterbox.md for one-time DNS + NPM steps.

DEPLOY_HOST  ?= liu
DEPLOY_WS    ?= /home/liu/infra/pointegrityws
DEPLOY_SSH   := ssh $(DEPLOY_HOST)

.PHONY: build-letterbox
build-letterbox: ## Build letterbox binary into ./build/
	go build -o ./build/letterbox ./cmd/letterbox

.PHONY: dev-letterbox
dev-letterbox: build-letterbox ## Run letterbox locally on :3737
	@mkdir -p data
	LETTERBOX_ADDR=:3737 \
	LETTERBOX_DB=./data/letterbox.db \
	LETTERBOX_REDIRECT=http://localhost:$(PORT) \
	LETTERBOX_ADMIN_KEY=devkey \
	./build/letterbox

.PHONY: deploy-pull-letterbox
deploy-pull-letterbox: ## git pull pointegrity on remote
	$(DEPLOY_SSH) "sudo -u liu bash -lc 'cd $(DEPLOY_WS)/pointegrity && git pull --ff-only'"

.PHONY: deploy-build-letterbox
deploy-build-letterbox: ## Rebuild letterbox on remote
	$(DEPLOY_SSH) "sudo -u liu bash -lc 'export PATH=/usr/local/go/bin:\$$PATH && cd $(DEPLOY_WS)/pointegrity && go build -o ./build/letterbox ./cmd/letterbox && ls -l ./build/letterbox'"

.PHONY: deploy-restart-letterbox
deploy-restart-letterbox: ## Restart letterbox systemd service
	$(DEPLOY_SSH) sudo systemctl restart letterbox

.PHONY: deploy-letterbox
deploy-letterbox: deploy-pull-letterbox deploy-build-letterbox deploy-restart-letterbox ## Pull + build + restart letterbox

.PHONY: deploy-status-letterbox
deploy-status-letterbox: ## Show letterbox service status
	$(DEPLOY_SSH) sudo systemctl status letterbox --no-pager

.PHONY: deploy-logs-letterbox
deploy-logs-letterbox: ## Tail letterbox logs
	$(DEPLOY_SSH) sudo journalctl -u letterbox -f --no-pager -n 50
