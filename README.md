TODO

- Add default params to the container, so the limit and AssumeRole could be changed via arguments

Example curls

curl 127.0.0.1:8080 -d '{test: "test"}'

curl 127.0.0.1:8080?message=test

curl '127.0.0.1:8080?message=test&subject=test'

curl '127.0.0.1:8080?subject=test' -d '{test: "test"}'
