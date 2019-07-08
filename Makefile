VERSION=$(shell git describe --always --long)
USER=stevelacy
NAME=dbjumper
IMAGE=$(USER)/$(NAME):$(VERSION)

all: docker

docker:
	docker build -t $(IMAGE) .
push:
	docker push $(IMAGE)
