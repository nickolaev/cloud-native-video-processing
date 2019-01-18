TOP := $(shell pwd)
DOCKER_DIR := $(TOP)/deployments/Docker
VERSION := latest
KUBERNETES_DIR := $(TOP)/deployments/Kubernetes
ISTIO_DIR := $(TOP)/deployments/Istio

TYPE := kubernetes
K8S_DEPLOYMENT := $(KUBERNETES_DIR)/cloud-native-video-process.yaml
ISTIO_DEPLOYMENT := $(ISTIO_DIR)/cloud-native-video-process.yaml
ISTIO_VIRTUAL_SERVICE_DEPLOYMENT := $(ISTIO_DIR)/cloud-native-video-process-istio.yaml

TARGETS := process source

default: deploy

.PHONY: build
build: $(addprefix video-,$(TARGETS))

video-%:
	docker build -t $*:$(VERSION) -f $(DOCKER_DIR)/Dockerfile.$* .

.PHONY: k8s-deploy
k8s-deploy:
	kubectl apply -f $(K8S_DEPLOYMENT)

.PHONY: k8s-delete
k8s-delete:
	kubectl delete -f $(K8S_DEPLOYMENT) || true

.PHONY: istio-deploy
istio-deploy:
	kubectl apply -f $(ISTIO_DEPLOYMENT)
	kubectl apply -f $(ISTIO_VIRTUAL_SERVICE_DEPLOYMENT)

.PHONY: istio-delete
istio-delete:
	kubectl delete -f $(ISTIO_VIRTUAL_SERVICE_DEPLOYMENT) || true
	kubectl delete -f $(ISTIO_DEPLOYMENT) || true

.PHONY: distclean
distclean: k8s-delete istio-delete
	docker system prune
