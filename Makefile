REST_MAIN := "$(CURDIR)/cmd/rest"
BIN_REST := "$(CURDIR)/bin/rest-service"
# EVENT_MAIN := "$(CURDIR)/cmd/event"
# BIN_EVENT := "$(CURDIR)/bin/event-service"


fetch:
	@go mod download

build-rest:
	@go build -i -v -o $(BIN_REST) $(REST_MAIN)

# build-event:
# 	@go build -i -v -o $(BIN_EVENT) $(EVENT_MAIN)

build-rest-vendor:
	@go build -mod=vendor -ldflags="-w -s" -o $(BIN_REST) $(REST_MAIN)

# build-event-vendor:
# 	@go build -mod=vendor -ldflags="-w -s" -o $(BIN_EVENT) $(EVENT_MAIN)

deploy: clean fetch build-rest 

swagger-gen:
	@swag init -g cmd/rest/main.go --output pkg/swagger/docs --generatedTime --md pkg/swagger/docs

clean:
	@rm -rf $(CURDIR)/bin

run-test:
	@go test -v -cover ./... | grep _test.go > output.out