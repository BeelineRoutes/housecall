
# go params
GOCMD=go

# normal entry points
	
update:
	clear 
	@echo "updating dependencies..."
	@go get -u -t ./...
	@go mod tidy 

build:
	clear 
	@echo "building..."
	@$(GOCMD) build .
	
test-first:
	clear
	@echo "testing housecall primary auth functions..."
	@$(GOCMD) test -run TestFirst ./...

test-second:
	clear
	@echo "test housecall second level functions..."
	@$(GOCMD) test -run TestSecond ./...

test-third:
	clear
	@echo "test housecall third level functions..."
	@$(GOCMD) test -run TestThird ./...
