
# go params
GOCMD=go
GOTEST=$(GOCMD) test

# normal entry points
	
update:
	@go get -u all

build:
	clear 
	@$(GOTEST) -run TestFirstModelsError ./...

test-first: build
test-first:
	clear
	@echo "testing housecall primary auth functions..."
	@$(GOTEST) -v -run TestFirst ./...

test-second: build
test-second:
	clear
	@echo "test housecall second level functions..."
	@$(GOTEST) -run TestSecond ./...

test-third: build
test-third:
	clear
	@echo "test housecall third level functions..."
	@$(GOTEST) -v -run TestThird ./...

