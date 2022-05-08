.PHONY: build run stop
build:
	docker build -t http-metric .

run:
	docker run -d --name http-metric -p 8080:80 http-metric

stop:
	docker kill http-metric

clean: stop
	docker system prune

start:
	docker start http-metric
