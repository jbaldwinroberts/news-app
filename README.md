# Esqimo
News app - backend API

A small API to serve RSS feeds.

## API documentation

View the API documentation [here](https://app.swaggerhub.com/apis-docs/josephroberts7/news-app_api/1.0.0).

## Godoc documentation
- [reader](https://godoc.org/github.com/josephroberts/news-app/reader)
- [server](https://godoc.org/github.com/josephroberts/news-app/server)
- [store](https://godoc.org/github.com/josephroberts/news-app/store)

## Prerequisites
This guide assumes you have git and docker installed on your local machine

## Build
```
$ git clone https://github.com/josephroberts/news-app.git
$ cd news-app
$ docker build -t news-app .
```

## Run
```
docker run -it -p 5050:5050 -v ${PWD}/config.yaml:/app/config.yaml  news-app
```

## Configure
Edit `config.yml` to add/remove feed sources and to change the port used - remember to change port in the docker run command too

## Assumptions and hacks
- Feed title is unique - if feed title is not unique, only one feed with that title will be served by the API
- When filtering by category, don't serve feed items that don't have a category
- In a feed item has more than one category I only consider the first when filtering by category

## Up next
- Write some tests
- Write some examples
- Replace `MemoryStore` with a DB or map `category` to feed item `GUID` and use this to filter by category more efficiently
- Read RSS feeds inside `go routines` to perform `GET` requests in parallel and handle timeouts
- Define `error` types and handle `errors` better