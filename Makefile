CONFIG_PATH := config/local.yml

run:
	go run cmd/url-shortener/main.go -CONFIG_PATH=$(CONFIG_PATH)
