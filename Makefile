dev:
	docker-compose build --no-cache scraper --build-arg mode=development
	docker-compose up --abort-on-container-exit

prod:
	docker-compose build --no-cache scraper --build-arg mode=production
	docker-compose up --abort-on-container-exit