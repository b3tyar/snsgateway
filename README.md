TODO

- Add default params to the container, so the limit and AssumeRole could be changed via arguments
- Add argument to favor the parameter over the POST data for the "message"

Example curls

curl 127.0.0.1:8080 -d '{test: "test"}'

curl 127.0.0.1:8080?message=test

curl '127.0.0.1:8080?message=test&subject=test'

curl '127.0.0.1:8080?subject=test' -d '{test: "test"}'
