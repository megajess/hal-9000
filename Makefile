run-server:
	@cd server && go run main.go

run-tests:
	@cd server && go test -v ./...

register:
	@cd util && bash testing.sh register

setOn:
	@cd util && bash testing.sh setOn

setOff:
	@cd util && bash testing.sh setOff
