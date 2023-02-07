SSL_PATH := /etc/letsencrypt/live/vps.osac.org.np
SERVE_PATH := $(shell which serve)

dev:
	@ npm install
	@ npm run dev

build:
	@ npm install
	@ npm run build

run: build
	@ $(SERVE_PATH) ./dist --single --debug --listen 443 --ssl-cert $(SSL_PATH)/fullchain.pem --ssl-key $(SSL_PATH)/privkey.pem --no-clipboard
