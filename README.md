## Kartini

This a package manager inspired by crux, kiss linux package manager.
so i build Kartini.

### Setup

First, set environment variable.
```sh
export KARTINI_ROOT="/test" # put your desired root system path
export KARTINI_PATH=""
export KARTINI_CACHE=""
```
Then,

``` 
./kartini help
```

## Make Package

This format following TOML.
```toml
name="kartini_base"
version="0.1"
depends=[
  # this contain depend with a package
]
sources=[
  # this contain source
  ["http://uwu/kartini.tar.xz"]
]
```

### License

[MIT](./LICENSE)
