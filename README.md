# infra
Golang implementation of the REST service for infrastructure.

This code use [gokit package](https://gokit.io) to implement a REST HTTP service.

## Install and Run
$ go get github.com/xinyu/infra/inventory

$ go run main.go -http.addr :8080
ts=2019-10-31T04:36:24.867797016Z caller=main.go:56 transport=HTTP addr=:8080


## API
### Host
$ curl -d '{"id":"1001","Name":"host1001"}' -H "Content-Type: application/json" -X POST http://localhost:8080/host/v1/hostinfo/

$ curl localhost:8080/host/v1/hostinfo/1001

$ curl -d '{"id":"1001","Name":"host1001-01"}' -H "Content-Type: application/json" -X PUT http://localhost:8080/host/v1/hostinfo/1001

$ curl -X DELETE localhost:8080/host/v1/hostinfo/1001

### Service
$ curl -d '{"id":"100001","Name":"testapp001", "HostID":"1001"}' -H "Content-Type: application/json" -X POST http://localhost:8080/service/v1/serviceinfo/

$ curl localhost:8080/service/v1/serviceinfo/100001

$ curl -d '{"id":"100001","Name":"testapp001-01", "HostID":"1001"}' -H "Content-Type: application/json" -X PUT http://localhost:8080/service/v1/serviceinfo/100001

$ curl -X DELETE localhost:8080/service/v1/serviceinfo/100001

 
