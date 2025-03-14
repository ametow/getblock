# GetBlock.io

## Overview
GetBlock.io is a blockchain node provider that allows developers to connect to various blockchain networks via API. This project provides a simple way to interact with GetBlock.io services using both a command-line interface (CLI) and a web API.

## Getting Started
To use GetBlock.io, you need to configure an API key. This can be done by setting the `$API_KEY` variable in the `makefile` or assigning it to the environment variable `GETBLOCK_KEY`.

### Prerequisites
- Make sure you have `make` installed.
- An active API key from [GetBlock.io](https://getblock.io/).

## Usage

### Command-Line Interface (CLI)
To run the CLI version, use the following commands:

```sh
make cli
make run-cli
```

### Web API Version
To launch the web API service, run:

```sh
make service
make run-service
```

Once the service is running, you can access the API at:

[http://localhost:8080/api/get](http://localhost:8080/api/get)

## Running with Docker
For convenience, you can run the application using Docker:

```sh
make container-cli
```

or

```sh
make container-service
```

## Contributing
Contributions are welcome! Please submit a pull request or open an issue if you find any bugs or have feature requests.

## License
This project is licensed under the MIT License.

## Support
For any issues, feel free to reach out or open an issue in the repository.

