.PHONY: build run stop
build:
	docker build -t http01 .

run:
	docker run -d --name http01 -p 8080:80 http01

stop:
	docker kill http01

clean: stop
	docker system prune

start:
	docker start http01
