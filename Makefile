build-lite:
	cd frontend && npm install && npm run build
	cd backend && CGO=0 go build -o foodhubber && mv foodhubber ../env/

build-db:
	- cd env/ && rm foodhubber.db
	cd env/ && sqlite3 foodhubber.db < ../data/structure.sql
	cd env/ && sqlite3 foodhubber.db < ../data/data.sql
	cd env/ && ls ../data/data-private.sql && sqlite3 foodhubber.db < ../data/data-private.sql

build:
	cd frontend && npm install && npm run build
	cd backend && CGO=0 go build -a -tags netgo,osusergo -ldflags '-w -extldflags "-static"' -trimpath -o foodhubber && mv foodhubber ../env/

zbuild:
	make build
	strip env/foodhubber
	upx --ultra-brute env/foodhubber

run-devel:
	cd frontend && npm install && npm run build
	cd backend && go run main.go --force-week 1 --db ../env/foodhubber.db

update:
	cd frontend && npm update
	cd backend && go get -u && go mod tidy

clean:
	rm -rf frontend/node_modules
	- rm env/foodhubber*