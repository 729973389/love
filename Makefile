.PHONY:run dlv build clean tar wuff
run :
	@./build/edgeDaemon
dlv:
	go build cmd/edgeDaemon/main.go
	sh dlv.sh
build:
	@go build -o build/edgeDaemon cmd/edgeDaemon/main.go
	@go build -o build/jsonCreater jsonCreater/main.go
	@go build -o installEdgeDaemon/jsonCreater jsonCreater/main.go
	@go build -o installEdgeDaemon/edgeDaemon cmd/edgeDaemon/main.go
clean:
	@rm build/server.json || echo "ERROR"
	@rm build/edgeDaemon||echo "ERROR"
	@rm build/jsonCreater ||echo "ERROR"
	@rm installEdgeDaemon/edgeDaemon ||echo "ERROR"
	@rm installEdgeDaemon/jsonCreater || echo "ERROR"
	@rm /wuff1996/nginx/source/installEdgeDaemon.* || echo "ERROR"
tar:
	@go build -o installEdgeDaemon/jsonCreater jsonCreater/main.go
	@go build -o installEdgeDaemon/edgeDaemon cmd/edgeDaemon/main.go
	@sudo tar zcvf /home/wuff/windows/windows/installEdgeDaemon.tar.gz installEdgeDaemon
wuff:
	@sudo chown -R wuff:wuff ../edgeDaemon