#!/bin/bash
#
# Show differences against facter4
#
# go get github.com/homeport/dyff/cmd/dyff
#

TEMP=$(mktemp -d)
trap 'rm -rf $TEMP' EXIT

./ufacter-linux-amd64 -yaml > "$TEMP/ufacter.yaml"
./ufacter-linux-amd64 -yaml -no-volatile > "$TEMP/ufacter-no-volatile.yaml"

dyff between "$TEMP/ufacter.yaml" "$TEMP/ufacter-no-volatile.yaml"
