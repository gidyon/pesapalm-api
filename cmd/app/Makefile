local_image := pesapalm
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

run:
	go build -v -o service -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn" && \
	./service -config-file=.env -static-dir=out

compile:
	@GOOS=linux CGO_ENABLED=1 go build -tags netgo -installsuffix netgo -v -o service -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn" .

docker_build:
	docker build -t $(local_image) .

docker_tag:
	@docker tag $(local_image) $(image)

docker_push:
	@docker push $(image)

build_image: docker_build docker_tag docker_push

build_service: compile docker_build docker_tag docker_push

build_service_skaffold: compile docker_build docker_tag docker_push_skaffold


