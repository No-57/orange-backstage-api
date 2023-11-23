.PHONY: doc

API_LAYOUT_DESIGN = ./infra/api

doc: doc-fmt backstage-doc frontstage-doc

backstage-doc:
	swag init --dir ./cmd/server/,$(API_LAYOUT_DESIGN),./app

frontstage-doc:
	swag init --dir ./docs/frontstage,$(API_LAYOUT_DESIGN) --output ./docs/frontstage --ot json

doc-fmt:
	swag fmt -d ./,./docs/frontstage
