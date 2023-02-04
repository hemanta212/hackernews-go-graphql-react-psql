SSL_PATH := /etc/letsencrypt/live/vps.osac.org.np
SERVE_PATH := $(shell which serve)

dev:
	@ npm start

build:
	@ npm run build

run: build
	@ sudo $(SERVE_PATH) ./build --single --debug --listen 443 --ssl-cert $(SSL_PATH)/fullchain.pem --ssl-key $(SSL_PATH)/privkey.pem --no-clipboard
