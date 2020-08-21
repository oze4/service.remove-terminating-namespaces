CN=CONTAINER_NAME

.PHONY: docker
docker:
	docker build --pull --rm -f "Dockerfile" -t oze4/service.remove-terminating-namespaces:latest "."
	docker push oze4/service.remove-terminating-namespaces:latest

.PHONY: kubernetes
kubernetes:
	kubectl delete -f deploy/rbac.yaml
	kubectl delete -f deploy/cronjob.yaml
	kubectl apply -f deploy/rbac.yaml
	kubectl apply -f deploy/cronjob.yaml

.PHONY: kubernetes-apply
kubernetes-apply:
	kubectl apply -f deploy/rbac.yaml
	kubectl apply -f deploy/cronjob.yaml

.PHONY: kubernetes-delete
kubernetes-delete:
	kubectl delete -f deploy/rbac.yaml
	kubectl delete -f deploy/cronjob.yaml

.PHONY: deploy
deploy: docker kubernetes