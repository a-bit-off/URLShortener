CONFIG_PATH := config/local.yml

all: tests run

run:
	go run cmd/url-shortener/main.go -CONFIG_PATH=$(CONFIG_PATH)

tests:
	go test internal/http-server/handlers/url/save/save_test.go
	go test internal/http-server/handlers/redirect/redirect_test.go
	go test internal/http-server/handlers/delete/delete_test.go






