SET GOPATH=%~dp0
go build -o ./bin/GameServer.exe ./src/GameServer.go 
go build -o ./bin/LoginServer.exe ./src/LoginServer.go 
go build -o ./bin/extractor.exe ./src/ExtractorMain.go 