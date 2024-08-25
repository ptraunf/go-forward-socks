# Go Forward SOCKS Proxy

E2E Test:
`curl -x socks5[h]://<proxy> <target>`

Example:
```sh
# Resolve domain locally, send IP to proxy
curl -x socks5://127.0.0.1:8080 http://httpbin.org/get

# Resolve domain at proxy
curl -x socks5h://127.0.0.1:8080 http://httpbin.org/get
#             ^ 
```


