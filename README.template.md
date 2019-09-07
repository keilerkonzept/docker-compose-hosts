# ${APP}

Work with docker-compose manifests for multiple hosts.

## Contents

- [Contents](#contents)
- [Get it](#get-it)
  - [Using `go get`](#using-go-get)
  - [Pre-built binary](#pre-built-binary)
- [Usage](#usage)
  - [Command-line interface syntax](#command-line-interface-syntax)
  - [`docker-compose-hosts.yml` syntax](#docker-compose-hostsyml-syntax)


## Get it

### Using `go get`

```sh
go get -u github.com/sgreben/${APP}
```

### Pre-built binary

[Download a binary](https://github.com/sgreben/${APP}/releases/latest) from the releases page or from the shell:

```sh
# Linux
curl -L https://github.com/sgreben/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_linux_x86_64.tar.gz | tar xz

# OS X
curl -L https://github.com/sgreben/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_osx_x86_64.tar.gz | tar xz

# Windows
curl -LO https://github.com/sgreben/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_windows_x86_64.zip
unzip ${APP}_${VERSION}_windows_x86_64.zip
```

## Usage

`${APP}` reads project definitions from a [`${APP}.yml` file](#docker-compose-hostsyml-syntax), and forwards all positional arguments to `docker-compose` for each of the projects.

### Command-line interface syntax

```text
${APP} [OPTIONS] -- [COMMAND [ARGS...]]
```

```text
${USAGE}
```

### `docker-compose-hosts.yml` syntax

Example file with all fields below; see also [example/docker-compose-hosts.yml](example/docker-compose-hosts.yml).

```yaml
version: '0.1'
projects:
  project-name-goes-here:
    compose_file: (path to docker-compose.yml)
    docker_host: (value for DOCKER_HOST, optional)
  another-project-name-goes-here:
    compose_file: (path to docker-compose.yml)
    docker_host: (value for DOCKER_HOST, optional)
```

All string fields except `version` support `$ENVVARS`; A literal `$` can be produced using the escape `$$`.
