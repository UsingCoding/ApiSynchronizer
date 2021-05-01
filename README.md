## Dev tools to sync api files between services

### Build App

1. In `.env` set up your Apis repo url
2. Run `make`
3. Then `apisynchronizer --help` to get know how tool works

### Usage

1. In your directory you should have `build.yaml` or set config by passing option `-f`
2. There you should describe your api dependencies. Example below
```yaml
api:
    apigateway: master
    contentservice: master
    userservice: master
```
3. Afterward provide `-o` option to set up output directory. Example of command
```shell
apisynchronizer sync -o api -f build.yaml
```
4. You apis repo should have similar structure
```
api/
    apigateway.proto
    contentservice.proto
    userservice.swagger
```
5. After run `apisynchronizer sync -o api` your api files will be copied to...
```
api/
    apigateway/
        apigateway.proto
    contentservice/
        contentservice.proto
    userservice/
        userservice.swagger
```