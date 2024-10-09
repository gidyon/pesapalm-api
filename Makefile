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

protoc_jasmin_connector.v1:
	@protoc -I=$(API_IN_PATH)/v1/jasmin -I=third_party --go-grpc_out=$(API_OUT_PATH)/connectors/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/connectors/v1 --go_opt=paths=source_relative connectors.proto
	@protoc -I=$(API_IN_PATH)/v1/jasmin -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/connectors/v1 connectors.proto
	@protoc -I=$(API_IN_PATH)/v1/jasmin -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) connectors.proto

protoc_routes.v1:
	@protoc -I=$(API_IN_PATH)/v1/jasmin -I=third_party --go-grpc_out=$(API_OUT_PATH)/routes/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/routes/v1 --go_opt=paths=source_relative routes.proto
	@protoc -I=$(API_IN_PATH)/v1/jasmin -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/routes/v1 routes.proto
	@protoc -I=$(API_IN_PATH)/v1/jasmin -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) routes.proto

protoc_jasmin.v1:
	@protoc -I=$(API_IN_PATH)/v1/jasmin -I=third_party --go-grpc_out=$(API_OUT_PATH)/jasmin/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/jasmin/v1 --go_opt=paths=source_relative jasmin.proto
	@protoc -I=$(API_IN_PATH)/v1/jasmin -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/jasmin/v1 jasmin.proto
	@protoc -I=$(API_IN_PATH)/v1/jasmin -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) jasmin.proto

protoc_blacklist.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/blacklist/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/blacklist/v1 --go_opt=paths=source_relative blacklist.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/blacklist/v1 blacklist.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) blacklist.proto

protoc_bulksmsrates.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/bulksmsrates/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/bulksmsrates/v1 --go_opt=paths=source_relative bulksmsrates.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/bulksmsrates/v1 bulksmsrates.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) bulksmsrates.proto

protoc_campaign.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/campaign/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/campaign/v1 --go_opt=paths=source_relative campaign.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/campaign/v1 campaign.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) campaign.proto

protoc_client.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/client/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/client/v1 --go_opt=paths=source_relative client.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/client/v1 client.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) client.proto

protoc_contact.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/contact/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/contact/v1 --go_opt=paths=source_relative contact.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/contact/v1 contact.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) contact.proto

protoc_contactgroup.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/contactgroup/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/contactgroup/v1 --go_opt=paths=source_relative contactgroup.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/contactgroup/v1 contactgroup.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) contactgroup.proto

protoc_developer.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/developer/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/developer/v1 --go_opt=paths=source_relative developer.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/developer/v1 developer.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) developer.proto

protoc_longrunning.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/longrunning/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/longrunning/v1 --go_opt=paths=source_relative longrunning.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/longrunning/v1 longrunning.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) longrunning.proto

protoc_message_template.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/messagetemplate/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/messagetemplate/v1 --go_opt=paths=source_relative messagetemplate.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/messagetemplate/v1 messagetemplate.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) messagetemplate.proto

protoc_payment.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/payment/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/payment/v1 --go_opt=paths=source_relative payment.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/payment/v1 payment.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) payment.proto

protoc_purchase.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/purchase/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/purchase/v1 --go_opt=paths=source_relative purchase.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/purchase/v1 purchase.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) purchase.proto

protoc_senderid.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/senderid/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/senderid/v1 --go_opt=paths=source_relative senderid.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/senderid/v1 senderid.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) senderid.proto

protoc_sms.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/sms/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/sms/v1 --go_opt=paths=source_relative sms.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/sms/v1 sms.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) sms.proto

protoc_notification.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/notification/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/notification/v1 --go_opt=paths=source_relative notification.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/notification/v1 notification.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) notification.proto

protoc_auditlog.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/auditlog/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/auditlog/v1 --go_opt=paths=source_relative auditlog.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/auditlog/v1 auditlog.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) auditlog.proto

protoc_stat.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/stat/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/stat/v1 --go_opt=paths=source_relative stat.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/stat/v1 stat.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) stat.proto

protoc_emailing.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/emailing/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/emailing/v1 --go_opt=paths=source_relative emailing.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/emailing/v1 emailing.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) emailing.proto

protoc_apicreds.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/apicreds/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/apicreds/v1 --go_opt=paths=source_relative apicreds.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/apicreds/v1 apicreds.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) apicreds.proto

protoc_senderiddoc.v1:
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --go-grpc_out=$(API_OUT_PATH)/senderiddoc/v1 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/senderiddoc/v1 --go_opt=paths=source_relative document.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/senderiddoc/v1 document.proto
	@protoc -I=$(API_IN_PATH)/v1 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) document.proto

protoc_bulksmsrates.v2:
	@protoc -I=$(API_IN_PATH)/v2 -I=third_party --go-grpc_out=$(API_OUT_PATH)/bulksmsrates/v2 --go-grpc_opt=paths=source_relative --go_out=$(API_OUT_PATH)/bulksmsrates/v2 --go_opt=paths=source_relative bulksmsrates.proto
	@protoc -I=$(API_IN_PATH)/v2 -I=third_party --grpc-gateway_out=logtostderr=true,paths=source_relative:$(API_OUT_PATH)/bulksmsrates/v2 bulksmsrates.proto
	@protoc -I=$(API_IN_PATH)/v2 -I=third_party --openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:$(OPEN_API_V2_OUT_PATH) bulksmsrates.proto


copy_documentation:
	@cp -r $(OPEN_API_V2_OUT_PATH) cmd/apis/apidoc/

protoc_all: protoc_jasmin_connector.v1 protoc_routes.v1 protoc_jasmin.v1 protoc_blacklist.v1 protoc_bulksmsrates.v1 protoc_campaign.v1 protoc_client.v1 protoc_contact.v1 protoc_contactgroup.v1 protoc_developer.v1 protoc_longrunning.v1 protoc_message_template.v1 protoc_payment.v1 protoc_purchase.v1 protoc_senderid.v1 protoc_sms.v1 protoc_notification.v1 protoc_auditlog.v1 protoc_stat.v1 protoc_emailing.v1 protoc_apicreds.v1 protoc_senderiddoc.v1 protoc_bulksmsrates.v2 copy_documentation

local_image :=onfon-sms-api
image := public.ecr.aws/q1f9b5m5/onfon-sms-api
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

get_auth:
	aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/q1f9b5m5

docker_build:
	docker build -t $(local_image) .

docker_tag:
	@docker tag $(local_image) $(image)

docker_push:
	@docker push $(image)

build_service: docker_build docker_tag get_auth docker_push

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
