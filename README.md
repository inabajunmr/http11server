# http11server

HTTP/1.1 server

## Run

```
$ go run main.go
$ curl localhost -H "Accept: application/xml" -v
*   Trying ::1...
* TCP_NODELAY set
* Connection failed
* connect to ::1 port 80 failed: Connection refused
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 80 (#0)
> GET / HTTP/1.1
> Host: localhost
> User-Agent: curl/7.64.1
> Accept: application/xml
> 
< HTTP/1.1 200 OK
< Date: Sat, 02 Oct 2021 01:45:15 GMT
< Vary: accept-encoding, accept
< Accept-Range: bytes
< Content-Length: 240
< 
<Echo>
 <method>GET</method>
 <request_target>/</request_target>
 <version>HTTP/1.1</version>
 <headers>HOST: localhost</headers>
 <headers>USER-AGENT: curl/7.64.1</headers>
 <headers>ACCEPT: application/xml</headers>
 <body></body>
* Connection #0 to host localhost left intact
</Echo>* Closing connection 0
```

## Test

```
$ make test
```


## Support

* Chunked Request(only gzip and identity)
* Keey-Alive and Connection header
* HEAD/OPTION
* Content-Type
* Range Request
* Accept(only application/json and application/xml)
* Accept-Encoding

## Unsupported

* Pipeline
* Catche/Conditional Request
* chunked response
* multi-line header(in message/http)
* parse request target
* TE header
* Trailer
* multipart
* Accept-Charset
* Accept-Language

