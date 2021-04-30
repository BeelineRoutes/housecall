
# go params
GOCMD=go

# normal entry points
	
update:
	@go get -u all

build:
	clear 
	@$(GOCMD) test -run TestModelsError ./...

test:
	clear
	@echo "building house call..."
	@$(GOCMD) test ./...

