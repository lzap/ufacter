#!/bin/bash
#
# Show differences against facter4
#
# go get github.com/homeport/dyff/cmd/dyff
#

TEMP=$(mktemp -d)
trap 'rm -rf $TEMP' EXIT

facter -y > "$TEMP/facter.yaml"
./ufacter-linux-amd64 -yaml > "$TEMP/ufacter.yaml"

dyff between "$TEMP/facter.yaml" "$TEMP/ufacter.yaml"
