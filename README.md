# ufacter (micro facter)

A loose implementation of Puppet Labs facter in golang. The main goal are platforms where there isn't possible or feasible to install Ruby, and performance to be kept around 10 times faster than Facter 3.X.

Only very basic set of facts from Facter 3.X are implemented:

* cpu
* memory
* network
* storage
* kernel
* operating system

Facter version is reported in `facterversion` as `3.0.0`, there's also additional `ufacter.version` fact available.

## Differences and limitations

There are some facts missing at the moment which are made on purpose:

* Fact tree `identity` (uid/gid etc).
* Ruby version (not relevant).

Some things are not yet implemented but planned:

* Primary interface is not yet reported via `primary` fact.
* Additionally, IPv6 primary interface will be reported too via `primary6` fact.
* Primary interface facts (`ipaddress`, `network` and similar).
* Processor "total speed": `processors.speed`.
* Some names (e.g. OS distribution names) are be reported differently (see https://github.com/shirou/gopsutil/blob/master/host/host_linux.go).
* Operating system major, minor and LSB info (full name, description) are missing.
* Some units contain decimal part ("123.00 bytes" instead of "123 bytes").
* FIPS and SELinux modes.

## Additional facts

* network link - interface names, types and relations (bonds, vlans, bridges)
* `primary` and `primary6` device name in `network`

## Requirements

- go v1.5 or newer is required

The code has been tested on Linux at the moment.

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

- report netlink route tables: https://github.com/vishvananda/netlink/blob/master/route.go#L35
- report "primary" network interface via https://github.com/jackpal/gateway (both IPv4 and IPv6 - /proc/net/ipv6_route)
- stable only (constant of fact type or Add method)
- maximum interfaces/mountpoints/devices limit option
- better error handling
- report facter and ufacter versions