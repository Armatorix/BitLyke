OPENAPI_PATH=./api/openapi-spec/api.yml
GENERATE_PATH=./pkg/model
compose = docker-compose -f deployments/docker/docker-compose.yml

.PHONY: run
run:
	${compose} up --build -d

.PHONY: run-db-only
run-db-only:
	${compose} up -d postgres

.PHONY: stop
stop:
	${compose} down

.PHONY: model-rebuild
model-rebuild: model-remove model-generate model-clean format

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

.PHONY: linter
lint:
	golangci-lint run

.PHONY: formatter
format:
	goimports -w **/*.go 

test-e2e: 
	$(MAKE) run
	./scripts/wait_for_it.sh http://localhost:8081/public/health-check
	go test ./test/e2e
	$(MAKE) stop