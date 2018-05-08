GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=hello
ZIP_NAME=handler.zip

build:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINARY_NAME)
	zip ./$(ZIP_NAME) $(BINARY_NAME)
clean:
	-$(GOCLEAN) && rm $(BINARY_NAME)
	-rm -rf $(ZIP_NAME)