# docker-compose-hosts

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
go get -u github.com/sgreben/docker-compose-hosts
```

### Pre-built binary

[Download a binary](https://github.com/sgreben/docker-compose-hosts/releases/latest) from the releases page or from the shell:

```sh
# Linux
curl -L https://github.com/sgreben/docker-compose-hosts/releases/download/0.1.1/docker-compose-hosts_0.1.1_linux_x86_64.tar.gz | tar xz

# OS X
curl -L https://github.com/sgreben/docker-compose-hosts/releases/download/0.1.1/docker-compose-hosts_0.1.1_osx_x86_64.tar.gz | tar xz

# Windows
curl -LO https://github.com/sgreben/docker-compose-hosts/releases/download/0.1.1/docker-compose-hosts_0.1.1_windows_x86_64.zip
unzip docker-compose-hosts_0.1.1_windows_x86_64.zip
```

## Usage

`docker-compose-hosts` reads project definitions from a [`docker-compose-hosts.yml` file](#docker-compose-hostsyml-syntax), and forwards all positional arguments to `docker-compose` for each of the projects.

### Command-line interface syntax

```text
docker-compose-hosts [OPTIONS] -- [COMMAND [ARGS...]]
```

```text
Usage of docker-compose-hosts:
  -f string
    	(alias for -file) (default "docker-compose-hosts.yml")
  -file string
    	specify an alternate compose-hosts file (default "docker-compose-hosts.yml")
  -logs-off
    	disable all logging
  -logs-verbose
    	enable extra logging
  -parallel
    	run commands in parallel
  -q	(alias for -logs-off)
  -v	(alias for -logs-verbose)
  -version
    	print version and exit
```

### `docker-compose-hosts.yml` syntax

Example file with all fields below; see also [example/docker-compose-hosts.yml](example/docker-compose-hosts.yml).

```yaml
version: '1.0'
projects:
  project-name-goes-here:
    compose_file: (path to docker-compose.yml)
    docker_host: (value for DOCKER_HOST, optional)
  another-project-name-goes-here:
    compose_file: (path to docker-compose.yml)
    docker_host: (value for DOCKER_HOST, optional)
```

All string fields except `version` support `$ENVVARS`; A literal `$` can be produced using the escape `$$`.
