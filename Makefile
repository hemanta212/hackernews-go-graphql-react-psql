.DEFAULT_GOAL := all

.PHONY: all
all: clean backend frontend

.PHONY: clean
clean:
	rm backend/tmp.yaml
	rm frontend/tmp.yaml

.PHONY: frontend
frontend:
	cd frontend && make minikube

.PHONY: backend
backend:
	cd backend && make minikube
