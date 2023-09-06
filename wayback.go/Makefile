
build-dev:
	go build -gcflags=all='-N -l'

build-release:
	rm wayback.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

test-dev:
	go run main.go -targets http://testphp.vulnweb.com
	go run main.go -targets http://testphp.vulnweb.com,https://ginandjuice.shop
	go run main.go -details -targets http://testphp.vulnweb.com/,https://ginandjuice.shop

test-release:
	./wayback.go -details -targets http://testphp.vulnweb.com/,https://ginandjuice.shop