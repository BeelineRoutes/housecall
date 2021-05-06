
# go params
GOCMD=go

# normal entry points
	
update:
	@go get -u all

build:
	clear 
	@$(GOCMD) test -run TestModelsError ./...

test-second: build
test-second:
	clear
	@echo "test housecall second level functions..."
	@$(GOCMD) test -run TestHouseCallSecond ./...

test-first: build
test-first:
	clear
	@echo "testing housecall primary auth functions..."
	@$(GOCMD) test -run TestHouseCallFirst ./...

