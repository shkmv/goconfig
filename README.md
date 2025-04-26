# goconfig

Simple, flexible, and secure configuration loader for Go.
Load configurations from files, environment variables, and safely bind them into structs — with built-in secrets masking.
Inspired by Rust crate `config`.

---

## Features

- [x] Load from YAML files and environment variables
- [x] Merge multiple sources with priority
- [x] Bind into strongly-typed structs using tags
- [x] Minimalistic, clean API
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
    var cfg Config
    err := goconfig.New().
        FromFile("server-config.yaml").
        FromEnv("APP_").
        Bind(&cfg)
    if err != nil {
        panic(err)
    }
}

_ = cfg
```

###  Generics

```go
package main

import (
    "fmt"

    "github.com/shkmv/goconfig"
)

type ServerConfig struct {
	Port   int      `config:"port"`
  DBHost string   `config:"db.host"`
  DBPass string   `config:"db.pass"`
}

func main() {
cfg, err := goconfig.Load[Config](
    goconfig.WithFile("server-config.yaml"),
    goconfig.WithEnv("APP_"),
)
if err != nil {
    panic(err)
}

_ = cfg
```
