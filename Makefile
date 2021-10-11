DOCKER_REPO = gusampaio/hbday
DOCKER_TAG = latest

test:
	. scripts/setupTest.sh

run:
	docker compose build
	docker compose up

build: test
	docker build -t $(DOCKER_REPO):$(DOCKER_TAG) .

push: build
	docker login
	docker push $(DOCKER_REPO):$(DOCKER_TAG)

deploy:
	helm install hbday backend-chart --values backend-chart/values.yaml ||  helm upgrade hbday backend-chart --values backend-chart/values.yaml
