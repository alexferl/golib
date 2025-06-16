# golib

Collection of small Go modules I use across projects.

## Modules

- **config** - Configuration management with [viper](https://github.com/spf13/viper) and [pflag](https://github.com/spf13/pflag)
- **logger** - Structured logging wrapper around [zerolog](https://github.com/rs/zerolog)
- **http** - HTTP server and middleware components using [Echo](https://github.com/labstack/echo)
    - **middleware** - Configurable middleware components with CLI flag support
    - **server** - HTTP/HTTPS server with advanced features

## Installation

```shell
go get github.com/alexferl/golib/config
go get github.com/alexferl/golib/logger
go get github.com/alexferl/golib/http/middleware
go get github.com/alexferl/golib/http/server
```

## Usage

See individual module directories for examples and documentation.
