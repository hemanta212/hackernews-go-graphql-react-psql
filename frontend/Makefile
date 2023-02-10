SERVE_PATH := $(shell which serve)

dev:
	@ npm install
	@ npm run dev

build:
	@ npm install
	@ npm run build

deploy:
	@ if [[ -z "${HTTPS_SSL_PATH}" || -z "${SSL_CERT_FILE}" || -z "${SSL_KEY_FILE}" ]]; then\
		$(SERVE_PATH) ./dist --single --debug --listen 443 --no-clipboard;\
	else\
		$(SERVE_PATH) ./dist --single --debug --listen 443 --ssl-cert $(HTTPS_SSL_PATH)/$(SSL_CERT_FILE) --ssl-key $(HTTPS_SSL_PATH)/$(SSL_KEY_FILE) --no-clipboard;\
	fi

run: build deploy
