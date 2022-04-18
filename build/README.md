## Docker Build

### Build your own Docker image

1. Clone Project
```
git clone https://github.com/armosec/http-request.git kubescape && cd "$_"
```

2. Build
```
docker build -t http_request -f build/Dockerfile .
```