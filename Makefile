.PHONY: build run watch test docker-up docker-down list progress hint reset verify clean

BINARY := dragonflylings
CMD := ./cmd/dragonflylings

build:
	go build -o $(BINARY) $(CMD)

run: build
	./$(BINARY) run

watch: build
	./$(BINARY) watch

list: build
	./$(BINARY) list

progress: build
	./$(BINARY) progress

hint: build
	./$(BINARY) hint

reset: build
	./$(BINARY) reset

verify: build
	./$(BINARY) verify

test:
	go test ./... -v -count=1 -timeout 60s

docker-up:
	docker compose up -d
	@echo "Dragonfly running on port 6380"

docker-down:
	docker compose down

clean:
	rm -f $(BINARY) .dragonflylings-state.json

setup: docker-up
	@sleep 2
	@echo "Setup complete. Run 'make list' to see exercises."
