# goconfig

Simple, flexible, and secure configuration loader for Go.
Load configurations from files, environment variables, and safely bind them into structs — with built-in secrets masking.
Inspired by Rust crate `config`.

---

## Features

- [x] Load from YAML files and environment variables
- [x] Load from .env files
- [x] Merge multiple sources with priority
- [x] Bind into strongly-typed structs using tags
- [x] Minimalistic, clean API
- [x] Mark required fields with `required:"true"`
- [ ] Mask sensitive fields for secure logging

---

## Installation

```bash
go get github.com/shkmv/goconfig
```

## Usage

### Builder

```go
package main

import (
    "fmt"

    "github.com/shkmv/goconfig"
)

type ServerConfig struct {
    Port   int    `config:"port"`
    DBHost string `config:"db.host"`
    DBPass string `config:"db.pass"`
}

func main() {
    var cfg ServerConfig
    err := goconfig.New().
        FromFile("server-config.yaml").
        FromDotEnv(".env").
        FromEnv("APP_").
        Bind(&cfg)
    if err != nil {
        panic(err)
    }
}

_ = cfg
```

### Required fields

Add `required:"true"` to any field with a `config` tag to fail binding when the key is missing in all sources.

```go
type ServerConfig struct {
    Port   int    `config:"port" required:"true"`
    DBHost string `config:"db.host" required:"true"`
    DBPass string `config:"db.pass"` // optional
}
```

###  Generics

```go
package main

import (
    "fmt"

    "github.com/shkmv/goconfig"
)

type ServerConfig struct {
    Port   int    `config:"port"`
    DBHost string `config:"db.host"`
    DBPass string `config:"db.pass"`
}

func main() {
cfg, err := goconfig.Load[ServerConfig](
    goconfig.WithFile("server-config.yaml"),
    goconfig.WithDotEnv(".env"),
    goconfig.WithEnv("APP_"),
)
if err != nil {
    panic(err)
}
    
_ = cfg
```
### .env file format

Simple KEY=VALUE lines are supported. Lines beginning with `#` are comments. Optional `export` is allowed. Inline comments after unescaped `#` are stripped. Quotes and a few escapes (\n, \t, \r, \\) are handled.

Keys are normalized like environment variables: underscores become dots and keys are lowercased. For example `DB_HOST=localhost` becomes `db.host`.
