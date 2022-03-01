.PHONY: build

build:
	go build -o bin/main ./cmd

.PHONY: docker

run:
	@go run ./cmd

docker:
	@echo "Building Docker Image:"
	docker image build -f Dockerfile -t forum-image .
	@echo

	@echo "List of images:"
	docker images
	@echo

	@echo "Initiating Container:"
	docker container run -t -p 8080:8080 --detach --name forum-container forum-image
	@echo

	@echo "Running command:"
	docker exec -it forum-container ls -la
	@echo

	@echo "Running server:"
	docker exec -it forum-container ./main
	@echo

.PHONY: clean

start:
	# run existing container
	@docker start forum-container

stop: 
	@echo stopping running container:
	@docker stop forum-container 
remove:
	@echo remove everything images, containers, and networks
	docker image prune


clean:
	@echo "Stopping container:"
	docker stop forum-container
	@echo

	@echo "Removing container:"
	docker rm forum-container
	@echo

	@echo "Deleting images:"
	docker rmi -f forum-image
	@echo

	@echo "List of images and containers now:"
	docker ps -a
	@echo
	docker images
	@echo

	rm -rf main

.DEFAULT_GOAL := build