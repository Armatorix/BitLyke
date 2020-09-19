OPENAPI_PATH=./api/openapi-spec/api.yml
GENERATE_PATH=./pkg/model

.PHONY: model-rebuild
model-rebuild: model-remove model-generate model-clean

.PHONY: model-generate
model-generate:
	docker run --rm \
  		-v ${PWD}:/local openapitools/openapi-generator-cli generate \
		  -i /local/${OPENAPI_PATH} \
		  -g go \
		  -o /local/${GENERATE_PATH}
	sudo chown -R ${USER}:${USER} ./${GENERATE_PATH}

.PHONY: model-remove
model-remove:
	-rm -rf pkg/model

.PHONY: model-clean
model-clean:
	cd ${GENERATE_PATH} && \
	rm -rf .openapi-generator \
		api \
		docs 

	cd ${GENERATE_PATH} && \
	rm -f git_push.sh \
		go.mod \
		go.sum \
		README.md \
		.travis.yml