
# go params
GOCMD=go

# normal entry points
	
update:
	@go get -u all

test:
	clear
	@echo "building house call..."
	@$(GOCMD) test ./...

