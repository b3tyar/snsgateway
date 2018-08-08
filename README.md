TODO
- ~~pass data as POST param~~
- ~~proper logging~~ 
- ~~sanity/perf test~~
- ~~add the params as env variables to the container~~ 
- create JSON mode to process JSON body
- make the body processing bit a module that can be changed as part of the config
- make the listening port an env variable
- make the limit for the number of message an env variable


Example curl

curl 127.0.0.1:8080?key='{"Records\":[{"myRecords":"record"}]}' -g
