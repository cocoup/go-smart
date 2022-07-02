debug:
	go build -gcflags "all=-N -l" -o gocli-debug gos-cli.go

build:
	go build -ldflags="-s -w" -o gocli gocli.go
	$(if $(shell command -v upx), upx gocli)

mac:
	GOOS=darwin go build -ldflags="-s -w" -o gocli-darwin gocli.go
	$(if $(shell command -v upx), upx gocli-darwin)

win:
	GOOS=windows go build -ldflags="-s -w" -o gocli.exe gocli.go
	$(if $(shell command -v upx), upx gocli.exe)

linux:
	GOOS=linux go build -ldflags="-s -w" -o gocli-linux gocli.go
	$(if $(shell command -v upx), upx gocli-linux)

#image:
#	docker build --rm --platform linux/amd64 -t go-smart/gocli:$(version) .
#	docker tag go-smart/gocli:$(version) go-smart/gocli:latest
#	docker push go-smart/gocli:$(version)
#	docker push go-smart/gocli:latest
#	docker build --rm --platform linux/arm64 -t go-smart/gocli:$(version)-arm64 .
#	docker tag go-smart/gocli:$(version)-arm64 go-smart/gocli:latest-arm64
#	docker push go-smart/gocli:$(version)-arm64
#	docker push go-smart/gocli:latest-arm64
