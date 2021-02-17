# tcli

![Go](https://github.com/l-lin/tcli/workflows/Go/badge.svg)

> An interactive Trello client with auto-completion feature

![tcli](./tcli.gif)

## Installation
### Downloading standalone binary

Binaries are available from [Github releases](https://github.com/l-lin/tcli/releases).

### Using cURL

```bash
curl -sf https://gobinaries.com/l-lin/tcli | sh;
```

### Using docker

```bash
docker run -it --rm -v /path/to/.tcli.yml:/.tcli.yml ghcr.io/l-lin/tcli
```

### Building from source

```bash
# Build
make compile
```

## Usage

```bash
# explore the CLI with the help command
tcli -h

# start interactive mode
tcli

# you can also use it as a CLI
tcli ls /
```
