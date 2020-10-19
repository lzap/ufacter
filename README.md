# ufacter (micro facter)

A loose implementation of Puppet Labs facter in golang. The main goal are platforms where there isn't possible or feasible to install Ruby, and performance to be kept around 10 times faster than Facter 3.X. It is important that ufacter does NOT aim to be compatible with the original facter in any way, it just reports the most important facts the same way.

Only very basic set of facts from Facter 3.X are implemented:

* cpu
* memory
* network
* storage
* kernel
* operating system

Facter version is reported in `facterversion` as `3.0.0`, there's also additional `ufacter.version` fact available.

## Features

* Lightweight and fast (zero processes spawned during execution).
* YAML (default) and JSON output (does not support Ruby output).

## Differences and limitations

There are some facts missing at the moment which are made on purpose:

* Fact tree `identity` (uid/gid etc).
* Ruby version (not relevant).

## Additional facts

* network link - interface names, types and relations (bonds, vlans, bridges)
* `primary` and `primary6` device name in `network`

## Requirements

- go v1.12 or newer is required

The code has been tested on Linux at the moment.

## Build and use

```
go get -v github.com/lzap/ufacter/cmd/ufacter
cd ~/go/src/github.com/lzap/ufacter
go build ./cmd/ufacter
./ufacter -help
```

## Examples

```
# ./ufacter
disks:
  vda:
    vendor: "0x1af4"
is_virtual: true
kernel: Linux
kernelmajversion: "4.18"
kernelrelease: 4.18.0-147.5.1.el8_1.x86_64
kernelversion: 4.18.0
link:
  enp1s0:
    mac: 52:54:00:aa:bb:cc
    type: device
  enp7s0:
    mac: 52:54:00:dd:ee:ff
    type: device
  lo:
    type: device
memory:
  swap:
    available: 615.00 MiB
    available_bytes: 644870144
    capacity: 0.00%
    total: 615.00 MiB
    total_bytes: 644870144
    used: 0.00 bytes
    used_bytes: 0
  system:
    available: 1.27 GiB
    available_bytes: 1367322624
    capacity: 10.99%
    total: 1.60 GiB
    total_bytes: 1712922624
    used: 179.59 MiB
    used_bytes: 188317696
mountpoints:
  /:
    available: 2.21 GiB
    available_bytes: 2373533696
    capacity: 49.61%
    device: /dev/vda4
    filesystem: xfs
    options:
    - rw
    - relatime
    size: 4.39 GiB
    size_bytes: 4710203392
    used: 2.18 GiB
    used_bytes: 2336669696
  /boot:
    available: 668.23 MiB
    available_bytes: 700690432
    capacity: 26.46%
    device: /dev/vda2
    filesystem: ext4
    options:
    - rw
    - relatime
    size: 975.90 MiB
    size_bytes: 1023303680
    used: 240.47 MiB
    used_bytes: 252149760
  vda:
    size: 6.00 GiB
    size_bytes: 6442450944
network:
  primary: enp1s0
networking:
  fqdn: client2
  hostname: client2
  interfaces:
    enp1s0:
      bindings:
      - address: 192.168.122.5
        cidr: 192.168.122.5/24
        netmask: 255.255.255.0
        network: 192.168.122.0
      bindings6:
      - address: fe80::d171:5116:9b73:c43a
        cidr: fe80::d171:5116:9b73:c43a/64
        netmask: ffff:ffff:ffff:ffff:00:00:00:00
        network: 'fe80::'
      mac: 52:54:00:aa:bb:cc
      mtu: 1500
    enp7s0:
      mac: 52:54:00:dd:ee:ff
      mtu: 1500
    lo:
      bindings:
      - address: 127.0.0.1
        cidr: 127.0.0.1/8
        netmask: 255.0.0.0
        network: 127.0.0.0
      bindings6:
      - address: ::1
        cidr: ::1/128
        netmask: ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff
        network: ::1
      mtu: 65536
os:
  architecture: x86_64
  family: rhel
  hardware: x86_64
  name: centos
  release:
    full: 8.1.1911
path: /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/opt/puppetlabs/bin:/root/bin
processors:
  count: 1
  models:
  - AMD EPYC Processor (with IBPB)
  physicalcount: 1
system_uptime:
  boot_time: 1585665220
  days: 0
  hours: 1
  seconds: 5283
  uptime: 0 days
timezone: EDT
ufacter:
  errors:
    'link: IPv6 route': 101
virtual: ""
```

More interesting example is for VLAN over bonded interfaces, let's display only the link fact tree:

```
$ ufacter -modules link -json
{
  "link": {
    "bond0": {
      "bond": {
        "mode": 1
      },
      "mac": "52:54:00:aa:bb:cc",
      "type": "bond"
    },
    "bond0.10": {
      "mac": "52:54:00:aa:bb:cc",
      "parent": "bond0",
      "type": "vlan",
      "vlan": {
        "id": 10,
        "protocol": 33024
      }
    },
    "enp1s0": {
      "mac": "52:54:00:aa:bb:cc",
      "master": "bond0",
      "slave": "bond",
      "type": "device"
    },
    "enp7s0": {
      "mac": "52:54:00:aa:bb:cc",
      "master": "bond0",
      "slave": "bond",
      "type": "device"
    },
    "lo": {
      "type": "device"
    }
  }
}
```

This tool was written specifically for this major use case - to be used from the original facter via custom facts for netlink low-level information about interfaces (interface types and relations):

```ruby
Facter.add(:link) do
  setcode do
    json = `ufacter -json -modules link`
    JSON.parse(json)['link']
  end
end
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

Planned facts:

* IPv6 primary interface will be reported too via `primary6` fact.
* Primary interface facts (`ipaddress`, `network` and similar).
* Processor "total speed": `processors.speed`.
* Some names (e.g. OS distribution names) are be reported differently (see https://github.com/shirou/gopsutil/blob/master/host/host_linux.go).
* Operating system major, minor and LSB info (full name, description) are missing.
* Fact caching (`-check-new-facts` option returns 1 if there are modified/new facts) for cronjobs to only send updates when necessary
* Some units contain decimal part ("123.00 bytes" instead of "123 bytes").
* Output in shell mode (`FACT_NAME="fact value"`) so shell-based TUI could take advantage of this.
* Implement https://github.com/safchain/ethtool facts: link, speed, duplex, port, autoneg, wol (Linux only)
* FIPS and SELinux modes.
* IPMI facts from legacy discovery.
* Report EFI or BIOS mode (https://github.com/jcpunk/puppet-efi/blob/master/lib/facter/efi.rb)
* Routing table

Planned features:

- report netlink route tables: https://github.com/vishvananda/netlink/blob/master/route.go#L35
- report "primary" network interface via https://github.com/jackpal/gateway (both IPv4 and IPv6 - /proc/net/ipv6_route)
- stable only (constant of fact type or Add method)
- maximum interfaces/mountpoints/devices limit option
- better error handling
- report facter and ufacter versions
