build:
	docker build --no-cache -t foozie .

run:build
	docker run foozie