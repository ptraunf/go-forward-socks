TODOs:
- types for request, reply
- state machine - connect -> negotiate -> exchange -> close
- tests

E2E Test:
```sh
curl -x <proxy> <target>
curl -x 127.0.0.1:8080 http://httpbin.org/get
```


