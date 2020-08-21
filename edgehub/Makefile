.PHONY:run dlv build clean dev install
install:
	sudo sh install.sh
dev:
	touch edgeHub.log
run:
	./build/edgeHub
dlv:
	go build cmd/main.go
	sh dlv.sh
build:
	go build -o build/edgeHub cmd/main.go
clean:
	@rm http.json
	@rm edgeHub
