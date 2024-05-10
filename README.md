# engadget-scraper

### Running the scraper

Ensure your docker daemon is running as the app runs on Docker containers.
```shell
# For "development" mode that uses just 2 URLs to scrape:
make dev
```

```shell
# For running the complete scraper:
make prod
```

### Potential Improvements

* Error Handling: Some errors deserve a retry, while others don't, and how to surface these.
