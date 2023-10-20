SERVE_PATH := $(shell which serve)
APP?=hackernews-frontend
RELEASE?=0.0.1


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

container: build
	docker compose -f ../docker-compose.prod.yml build

minikube: container
	minikube image load $(APP):$(RELEASE)
	for t in $(shell find ../kubernetes/frontend -type f -name "*.yaml"); do \
	cat $$t | \
    	sed -E "s/\{ \{(\s*)\.Release(\s*)\} \}/$(RELEASE)/g" | \
    	sed -E "s/\{ \{(\s*)\.ServiceName(\s*)\} \}/$(APP)/g"; \
	echo ---; \
    done > tmp.yaml

	minikube kubectl -- apply -f tmp.yaml