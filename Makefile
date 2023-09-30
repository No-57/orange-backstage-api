.PHONY: doc

doc:
	swag fmt;swag init --dir ./cmd/server/,./app,./infra/api
