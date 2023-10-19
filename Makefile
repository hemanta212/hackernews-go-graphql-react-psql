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

.PHONY: portf
portf:
	minikube kubectl port-forward service/hackernews 8080:8080 &
	minikube kubectl port-forward service/hackernews-frontend 8000:80 &