name = imdb
img = golang:1.11
src = github.com/iuryfukuda/$(name)
workdir = /go/src/$(src)
run = docker run -v $(PWD):$(workdir) -w $(workdir) -it --rm $(img)

build:
	$(run) go build .

shell: image
	$(run) sh

run:
	$(run) go run main.go $(args)

check:
	$(run) go test ./...

check-args:
	$(run) go test $(args)


