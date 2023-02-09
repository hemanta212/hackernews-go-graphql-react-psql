SERVE_PATH := $(shell which serve)

dev:
	@ npm install
	@ npm run dev

build:
	@ npm install
	@ npm run build

deploy:
	@ if [[ -z $SSL_CERT || -z $SSL_KEY ]]; then\
		$(SERVE_PATH) ./dist --single --debug --listen 443 --no-clipboard;\
	else\
		$(SERVE_PATH) ./dist --single --debug --listen 443 --ssl-cert $(SSL_CERT) --ssl-key $(SSL_KEY) --no-clipboard;\
	fi

run: build deploy
