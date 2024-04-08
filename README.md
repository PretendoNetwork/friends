# Friends replacement server
Includes both the authentication and secure servers

## Compiling

### Setup
Install [Go](https://go.dev/doc/install) and [git](https://git-scm.com/downloads), then clone and enter the repository

```bash
$ git clone https://github.com/PretendoNetwork/friends
$ cd friends
```

### Compiling and running using `docker` (Preferred)
Install Docker either through your systems package manager or the [official installer](https://docs.docker.com/get-docker/)

To build the container:

```bash
$ docker build -t friends .
$ docker image prune --filter label=stage=builder -f
```
Optionally you may provide `BUILD_STRING` to `--build-arg` to set the authentication server build string

```bash
$ docker build -t friends --build-arg BUILD_STRING=auth-build-string .
$ docker image prune --filter label=stage=builder -f
```
If `BUILD_STRING` is not set, the default build string `pretendo.friends.docker` is used. You may also use the `docker` rule when building with `make` to set the build string automatically. See [compiling using `make`](#compiling-using-make) below for more info

To run the image first create a `.env` file with your [Configuration](#configuration) set before using `docker run`

Example:
```
PN_FRIENDS_POSTGRES_URI=postgres://username:password@localhost/friends?sslmode=disable
PN_FRIENDS_AUTHENTICATION_SERVER_PORT=60000
...
```

```bash
$ docker run --name friends --env-file .env -it friends
```

The image is compatible popular container managers such as Docker Compose and Portainer

### Compiling using `go`
To compile using Go, `go get` the required modules and then `go build` to your desired location. You may also want to tidy the go modules, though this is optional

```bash
$ go get -u
$ go mod tidy
$ go build -o build/friends
```

The server is now built to `build/friends`

When compiling with only Go, the authentication servers build string is not automatically set. This should not cause any issues with gameplay, but it means that the server build will not be visible in any packet dumps or logs a title may produce

To compile the servers with the authentication server build string, add `-ldflags "-X 'main.serverBuildString=BUILD_STRING_HERE'"` to the build command, or use `make` to compile the server

### Compiling using `make`
Compiling using `make` will read the local `.git` directory to create a dynamic authentication server build string, based on your repositories remote origin and current commit

Install `make` either through your systems package manager or the [official download](https://www.gnu.org/software/make/). We provide two different rules; A `default` rule which compiles [using `go`](#compiling-using-go), and a `docker` rule which compiles [using `docker`](#compiling-and-running-using-docker-preferred). Please refer to each sections setup instructions before continuing with your preferred rule

To build using `go`

```bash
$ make
```

The server is now built to `build/friends`

To build using `docker`

```bash
$ make docker
```

The image is now ready to run

## Configuration
All configuration options are handled via environment variables

`.env` files are supported

| Name                                        | Description                                                                                                            | Required                            |
|---------------------------------------------|------------------------------------------------------------------------------------------------------------------------|-------------------------------------|
| `PN_FRIENDS_CONFIG_DATABASE_URI`            | Fully qualified URI to your Postgres server (Example `postgres://username:password@localhost/friends?sslmode=disable`) | Yes                                 |
| `PN_FRIENDS_CONFIG_AUTHENTICATION_PASSWORD` | The password of the authentication server user account.                                                                | Yes                                 |
| `PN_FRIENDS_CONFIG_SECURE_PASSWORD`         | The password of the secure server user account. Used as part of the internal server data in Kerberos tickets           | Yes                                 |
| `PN_FRIENDS_CONFIG_AES_KEY`                 | AES key used in tokens provided by the account server                                                                  | Yes                                 |
| `PN_FRIENDS_CONFIG_GRPC_API_KEY`            | API key for your GRPC server                                                                                           | No (Assumed to be an open gRPC API) |
| `PN_FRIENDS_GRPC_SERVER_PORT`               | Port for the GRPC server                                                                                               | Yes                                 |
| `PN_FRIENDS_AUTHENTICATION_SERVER_PORT`     | Port for the authentication server                                                                                     | Yes                                 |
| `PN_FRIENDS_SECURE_SERVER_HOST`             | Host name for the secure server (should point to the same address as the authentication server)                        | Yes                                 |
| `PN_FRIENDS_SECURE_SERVER_PORT`             | Port for the secure server                                                                                             | Yes                                 |
| `PN_FRIENDS_ACCOUNT_GRPC_HOST`              | Host name for your account server gRPC service                                                                         | Yes                                 |
| `PN_FRIENDS_ACCOUNT_GRPC_PORT`              | Port for your account server gRPC service                                                                              | Yes                                 |
| `PN_FRIENDS_ACCOUNT_GRPC_API_KEY`           | API key for your account server gRPC service                                                                           | No (Assumed to be an open gRPC API) |
