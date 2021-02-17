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

## Inspiration

I'm mostly using the command line, and I like using VIM to edit contents. The Trello web UI is great, but I prefer
staying in the terminal. I could not find some good CLI / prompt, hence this project was born.

tcli was inspired by [trelew](https://github.com/fiatjaf/trelew) and unix commands, for their APIs are quite neat and
really powerful. Although the meaning of the commands do not really reflect the actions on the Trello resources, it's
still quite similar if we consider boards and lists as directories, and cards as files, thus avoiding the burden of
learning new commands.
