
Build
-----

Download the code

```
go get github.com/zenazn/goji/web
go get github.com/mdevilliers/take-home

```

Change directory to GOPATH/srcgithub.com/mdevilliers/take-home

Run the ci script - remember to set the execute flag if on linux

```
build\ci.sh
```

View the code coverage report

```
go tool cover -html=profile.cov
```

Run the app 

```
.\server
```

The rest server should be available from http://localhost:8000


Testing via curl
----------------

curl -I -X POST localhost:8000/topic1/user1

curl -X POST --data "message1" localhost:8000/topic1

curl -X POST --data "message2" localhost:8000/topic1

curl -v localhost:8000/topic1/user1

curl -I -X DELETE localhost:8000/topic1/user1
