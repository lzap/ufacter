#!/bin/bash
#
# Show differences against facter4
#
# go get github.com/homeport/dyff/cmd/dyff
#

TEMP=$(mktemp -d)
trap 'rm -rf $TEMP' EXIT

./ufacter -yaml > "$TEMP/ufacter.yaml"
./ufacter -yaml -no-extended > "$TEMP/ufacter-no-extended.yaml"

dyff between "$TEMP/ufacter.yaml" "$TEMP/ufacter-no-extended.yaml"
