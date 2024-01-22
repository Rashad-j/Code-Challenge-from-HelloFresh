unzip:
	$(info ******************** Unzipping files *********************************)
	@tar -xzvf files/hf_test_calculation_fixtures.tar.gz -C files
	@mv files/hf_test_calculation_fixtures.json files/fixtures.json

build:
	$(info ******************** Building parser *********************************)
	@go fmt ./...
	@go build -o bin/parser main.go

run: build
	$(info ******************** Running parser **********************************)
	@./bin/parser

test:
	$(info ******************** Running unit tests ******************************)
	@go test -v -race ./... -count=1

runWithoutStderr: build
	$(info ******************** Running parser without stderr *******************)
	@./bin/parser 2>/dev/null

runWithStderr: build
	$(info ******************** Running parser with stderr **********************)
	@./bin/parser stats --file ./files/test.json --postcode 10120 --words Potato,Mushroom,Veggie --fromTime 10AM --toTime 3PM

dockerBuildRunWithDefaultArgs:
	$(info ******************** Building and running parser in docker ************)
	@docker build -t parser-stats . && \
	docker run --rm -it parser-stats parser stats

dockerBuildRunWithCustomArgs:
	$(info ******************** Building and running parser in docker ************)
	@docker build -t parser-stats . && \
	docker run --rm -it parser-stats parser stats --file ./files/test.json --postcode 10120 --words Potato,Mushroom,Veggie --fromTime 9AM --toTime 2PM