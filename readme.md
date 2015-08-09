
Build
-----

Download the code

```
go get golang.org/x/tools/cmd/goimports
go get github.com/zenazn/goji/web
go get github.com/mdevilliers/take-home
```

Run the ci script - remember and set the execute flag on linux

```
build\ci.sh
```

View the code coverage report

```
go tool cover -html=profile.cov
```

Run the app 

```
go build github.com/mdevilliers/take-home/cmd/server/ 
.\server
```

The rest server should be available from http://localost:8080


Testing via curl
----------------

curl -I -X POST localhost:8000/topic1/user1

curl -X POST --data "message1" localhost:8000/topic1

curl -X POST --data "message2" localhost:8000/topic1

curl -v localhost:8000/topic1/user1

curl -I -X DELETE localhost:8000/topic1/user1
