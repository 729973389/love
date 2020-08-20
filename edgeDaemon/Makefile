.PHONY:run dlv build clean dev
dev:
	@touch edgeDaemon.log
run :
	@./edgeDaemon > edgeDaemon.log
dlv:
	go build cmd/main.go
	sh dlv.sh
build:
	go build -o build/edgeDaemon cmd/main.go
clean:
	@rm server.json
	@rm edgeDaemon
