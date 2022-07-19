debug:
	go build -gcflags "all=-N -l" -o {{.serviceName}}-debug {{.serviceName}}.go

build:
	go build -ldflags="-s -w" -o {{.serviceName}} {{.serviceName}}.go
	$(if $(shell command -v upx), upx {{.serviceName}})

mac:
	GOOS=darwin go build -ldflags="-s -w" -o {{.serviceName}}-darwin {{.serviceName}}.go
	$(if $(shell command -v upx), upx {{.serviceName}}-darwin)

win:
	GOOS=windows go build -ldflags="-s -w" -o {{.serviceName}}.exe {{.serviceName}}.go
	$(if $(shell command -v upx), upx {{.serviceName}}.exe)

linux:
	GOOS=linux go build -ldflags="-s -w" -o {{.serviceName}}-linux {{.serviceName}}d.go
	$(if $(shell command -v upx), upx {{.serviceName}}-linux)