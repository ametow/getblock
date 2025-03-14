# GetBlock.io  

## Launch  

To connect to the GetBlock.io API, you need to specify the API key in the `makefile` by setting the `$API_KEY` variable or assigning it to the environment variable `GETBLOCK_KEY`.  

### Command Line Version  

```sh
make cli
make run-cli
```

### Web API Version  

```sh
make service
make run-service
```

For the web version, the default endpoint is: [http://localhost:8080/api/get](http://localhost:8080/api/get)  

## Docker  

```sh
make container-cli
```

or  

```sh
make container-service
```
