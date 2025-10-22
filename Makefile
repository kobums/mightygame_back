tag=latest

all: server

server: dummy
	buildtool-model ./ 
	buildtool-router ./ > ./router/router.go
	go build -o bin/mighty main.go

fswatch:
	fswatch -0 controllers | xargs -0 -n1 build/notify.sh

run:
	gin --port 80 -a 9003 --bin bin/mighty run main.go

allrun:
	fswatch -0 controllers | xargs -0 -n1 build/notify.sh &
	gin --port 80 -a 9003 --bin bin/mighty run main.go

test: dummy
	go test -v ./...

linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/mighty.linux main.go

dockerbuild:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o bin/mighty.linux main.go

docker: dockerbuild
	docker build -t yiyuhki/mighty:$(tag) .

dockerrun:
	docker run -d --name="mighty" -p 9003:9003 yiyuhki/mighty

push: docker
	docker push yiyuhki/mighty:$(tag)

clean:
	rm -f bin/mighty

dummy:
