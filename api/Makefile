templ-install:
	@if ! command -v templ > /dev/null; then \
		read -p "Go's 'templ' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/a-h/templ/cmd/templ@latest; \
			if [ ! -x "$$(command -v templ)" ]; then \
				echo "templ installation failed. Exiting..."; \
				exit 1; \
			fi; \
		else \
			echo "You chose not to install templ. Exiting..."; \
			exit 1; \
		fi; \
	fi

build: templ-install
	@echo "Building..."
	@npm run build
	@templ generate
	@go generate
	@go build -o main main.go

up: 
	@docker compose -f ../local/docker-compose.yml up -d 

down:
	@docker compose -f ../local/docker-compose.yml down

test:
	@echo "Testing..."
	@STORAGE=postgres go test -v

.PHONY: templ-install build test up down