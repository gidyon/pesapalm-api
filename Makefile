API_IN_PATH := api/proto
API_OUT_PATH := pkg/api
OPEN_API_V2_OUT_PATH := api/openapiv2

setup_dev: ## Sets up a development environment for the okoa float bank apis project
	@cd deployments/dev &&\
	docker-compose up -d

setup_redis:
	@cd deployments/dev &&\
	docker-compose up -d redis

teardown_dev: ## Tear down development environment for the okoa float bank apis project
	@cd deployments/dev &&\
	docker-compose down

local_image := pesapalm-api
image := gidyon/pesapalm-api
context := .

ifdef IMAGE
	image=$(IMAGE)
else
	imagex := $(image)
	image_local := $(local_image)
	ifdef tag
		image=$(imagex):$(tag)
		local_image=$(image_local):$(tag)
	else	
		image=$(imagex):latest
		local_image=$(image_local)
	endif
endif

ifdef BUILD_CONTEXT
	context=$(BUILD_CONTEXT)
endif

docker_build:
	docker build -t $(local_image) .

docker_tag:
	@docker tag $(local_image) $(image)

docker_push:
	@docker push $(image)

build_service: docker_build docker_tag docker_push

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
