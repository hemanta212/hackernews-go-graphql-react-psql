SERVE_PATH := $(shell which serve)

dev:
	@ npm install
	@ npm run dev

build:
	@ npm install
	@ npm run build

deploy:
	@ if [[ -z "${HTTPS_SSL}" ]]; then\
		$(SERVE_PATH) ./dist --single --debug --listen 443 --no-clipboard;\
	else\
		$(SERVE_PATH) ./dist --single --debug --listen 443 --ssl-cert $(HTTPS_SSL)/fullchain.pem --ssl-key $(HTTPS_SSL)/privkey.pem --no-clipboard;\
	fi

run: build deploy
