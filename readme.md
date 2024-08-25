# Go Forward SOCKS Proxy

E2E Test:
`curl -x socks5h://<proxy> target`

Example:
```sh
# Resolve domain to IP locally 
curl -x socks5://127.0.0.1:8080 http://httpbin.org/get
```


