build-lite:
	cd frontend && npm install && npm run build
	cd backend && CGO=0 go build -o foodhubber && mv foodhubber ../

build:
	cd frontend && npm install && npm run build
	cd backend && CGO=0 go build -a -tags netgo,osusergo -ldflags '-w -extldflags "-static"' -trimpath -o foodhubber && mv foodhubber ../

zbuild:
	make build
	strip foodhubber
	upx --ultra-brute foodhubber

run-devel:
	cd frontend && npm install && npm run build
	cd backend && go run main.go --force-week 1

update:
	cd frontend && npm update
	cd backend && go get -u && go mod tidy