# ufacter (micro facter)

A loose implementation of Puppet Labs facter in golang. The main target are platforms where there isn't possible or feasible to install Ruby. Only very basic set of facts are implemented plus few more:

* cpu
* memory
* network (link layer - this is not provided by the original facter)
* network (IPv4 and IPv6)

## Requirements

- go v1.5 or newer is required

## Build and use

```
go get -v github.com/lzap/ufacter/cmd/ufacter
cd ~/go/src/github.com/lzap/ufacter
go build ./cmd/ufacter
./ufacter -help
```

## Environment variables

* `HOST_ETC` - specify alternative path to `/etc` directory
* `HOST_PROC` - specify alternative path to `/proc` mountpoint
* `HOST_SYS` - specify alternative path to `/sys` mountpoint

## Original work

Based on original work published at https://github.com/zstyblik/go-facter

## Licence

BSD 3-Clause ("BSD New" or "BSD Simplified") licence.

## TODO

- tree without dots (e.g. link.virbt-ntln.98.vlanprotocol)
- JSON output
- stable only (constant of fact type or Add method)
- maximum subtree
- calculate facts via goroutines

