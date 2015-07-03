all: go docker

go: 
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lazlo .

docker: 
	docker build -t lazlo .
