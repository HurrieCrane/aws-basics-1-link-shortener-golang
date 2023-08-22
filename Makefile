.PHONY: package

IMAGE_NAME=$(shell basename ${PWD})
IMAGE_TAG=${IMAGE_NAME}:local

ADDITIONAL_ARGS=

create_local_table:
	aws dynamodb create-table --table-name shortened-links --attribute-definitions AttributeName=link-hash,AttributeType=S --key-schema AttributeName=link-hash,KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 --endpoint http://localhost:8000

clean:
	rm -f $(wildcard lambdas/*/main)
	rm -f $(wildcard lambdas/*/main.zip)

test:
# count=1 forces it to not cache tests
	go test -count=1 ./pkg/...

package_local:
	$(MAKE) package -B ADDITIONAL_ARGS="-tags=local"

# Invoke generate like this:
# curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"queryStringParameters": { "uri": "https://www.google.com" } }'

# Invoke resolve like this:
# curl -XPOST "http://localhost:9001/2015-03-31/functions/function/invocations" -d '{ "pathParameters": { "hash": "68747470733a2f2f7777772e676f6f676c652e636f6dda39a3ee5e6b4b0d3255bfef95601890afd80709" } }'

run_local: package_local
	ENVIRONMENT=${ENV} AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} AWS_SESSION_TOKEN=${AWS_SESSION_TOKEN} docker compose up --build

# build and package each command
HANDLERS=$(addsuffix main,$(wildcard lambdas/*/))
$(HANDLERS): lambdas/%/*: *.go lambdas
	GOOS=linux
	GOARCH=amd64
	cd ./$(dir $@) && go build ${ADDITIONAL_ARGS} -o main . && zip -r -j main.zip main

package: $(HANDLERS)