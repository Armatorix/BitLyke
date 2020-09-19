OPENAPI_PATH=./api/openapi-spec/api.yml
GENERATE_PATH=./pkg/model


.PHONE: run
run:
	docker-compose \
			-f deployments/docker-compose.yml \
		up \
			--build

.PHONY: model-rebuild
model-rebuild: model-remove model-generate model-clean formatter

.PHONY: model-generate
model-generate:
	docker run --rm \
  			-v ${PWD}:/local openapitools/openapi-generator-cli \
		generate \
			-i /local/${OPENAPI_PATH} \
			-g go \
			-o /local/${GENERATE_PATH} \
			--package-name model

	sudo chown \
		-R ${USER}:${USER} \
		./${GENERATE_PATH}

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
		.travis.yml \
		.openapi-generator-ignore

.PHONY: prettier
prettier:
	docker run --rm \
			-v "$(pwd):" \
		tmknom/prettier \
			--parser=go \
			--write '**/*.go'

.PHONY: linter
linter:
	golangci-lint run

.PHONY: formatter
formatter:
	find . \
		-type f \
		-name \
		"*.go" \
		| xargs goimports -w