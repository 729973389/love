.PHONY:run dlv build clean dev
log:
	@touch edgeDaemon.log
	@./build/edgeDaemon > edgeDaemon.log
dlv:
	@go build -o build/edgeDaemon cmd/main.go
	@sh dlv.sh
build:
	@go build -o build/edgeDaemon cmd/main.go
	@go build -o build/jsonCreater jsonCreater/main.go
clean:
	@rm build/edgeDaemon
	@rm biuld/jsonCreater
