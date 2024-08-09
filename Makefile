GO = go
MAIN_FILE = main.go

default: serve

serve:
	$(GO) run $(MAIN_FILE) serve --http 0.0.0.0:8090

migrate:
	$(GO) run $(MAIN_FILE) migrate

dev:
	rm -fr pb_data; $(GO) run ${MAIN_FILE} migrate; $(GO) run $(MAIN_FILE) serve