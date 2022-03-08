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
	docker volume create sqlite3
	@echo

	@echo "Initiating Container:"
	docker run --rm -v sqlite3:app -p 27960:27960 --detach --name forum-container forum-image
	@echo

	@echo "Running command:"
	docker exec -it forum-container ls -la
	@echo

	@echo ---> Starting server on :27960 on: https://localhost:27960

.PHONY: clean

start:
	# run existing container
	docker run --rm -v sqlite3:/data -p 8080:8080 --detach --name forum-container forum-image
	@echo ---> Starting server on :8080 on: http://localhost:8080

stop: 
	@echo stopping running container:
	@docker stop forum-container 

clean:
	@echo Cleaning...
	@docker rmi forum-image
	@echo
	@echo docker container list:
	@docker container ls
	@echo
	@echo docker images list:
	@docker images

# remove everything from docker
rme:
	@docker system prune -a

.DEFAULT_GOAL := build